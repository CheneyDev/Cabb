package handlers

import (
    "context"
    "fmt"
    "sort"
    "strings"
    "sync"
    "time"

    "plane-integration/internal/cnb"
)

// acquireCreateLock locks on a (repo|planeIssueID) key to serialize CNB issue creation in-process
func (h *Handler) acquireCreateLock(repo, planeIssueID string) func() {
    if h == nil { return func() {} }
    key := repo + "|" + planeIssueID
    v, _ := h.createLocks.LoadOrStore(key, &sync.Mutex{})
    mu := v.(*sync.Mutex)
    mu.Lock()
    return func() { mu.Unlock() }
}

func (h *Handler) handlePlaneIssueEvent(env planeWebhookEnvelope, deliveryID string) {
    if !hHasDB(h) { return }
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    // Extract
    action := strings.ToLower(env.Action)
    data := env.Data
    planeIssueID, _ := dataGetString(data, "id")
    planeProjectID, _ := dataGetString(data, "project")
    name, _ := dataGetString(data, "name")
    descHTML, _ := dataGetString(data, "description_html")
    // Extract state name/id if present (for close decision)
    stateName, _ := dataGetStateNameAndID(data)
    // Fallback via activity for state changes
    if stateName == "" {
        if strings.EqualFold(strings.TrimSpace(env.Activity.Field), "state") {
            if s, ok := anyToStateName(env.Activity.NewValue); ok {
                stateName = s
            }
        }
    }
    // Dates
    startDate, dueDate := dataGetDates(data)
    // Priority (with DB mapping override when available)
    cnbPriority, planePriority, hasPriority := dataGetPriority(data)
    if hasPriority && planePriority != "" && planeProjectID != "" && hHasDB(h) {
        if mapped, ok, err := h.db.MapPlanePriorityToCNB(ctx, planeProjectID, planePriority); err == nil && ok {
            cnbPriority = mapped
        }
    }
    if !hasPriority {
        f := strings.ToLower(strings.TrimSpace(env.Activity.Field))
        if f == "priority" || f == "priority_name" || f == "priority_level" || f == "priority_value" {
            alt := map[string]any{"priority": env.Activity.NewValue}
            if cp, pp, ok := dataGetPriority(alt); ok {
                cnbPriority, planePriority, hasPriority = cp, pp, ok
            }
        }
    }
    assigneePlaneIDs := dataGetAssigneeIDs(data)

    labels := dataGetLabels(data)
    mappings, err := h.db.ListRepoProjectMappingsByPlaneProject(ctx, planeProjectID)

    LogStructured("info", map[string]any{
        "event":            "plane.issue",
        "delivery_id":      deliveryID,
        "action":           action,
        "plane_issue_id":   planeIssueID,
        "plane_project_id": planeProjectID,
        "labels":           labels,
        "start_date_set":   startDate != "",
        "due_date_set":     dueDate != "",
        "priority_set":     hasPriority,
        "plane_priority":   planePriority,
        "cnb_priority":     cnbPriority,
        "assignees_plane":  len(assigneePlaneIDs),
        "mappings_count":   len(mappings),
        "outbound_enabled": h.cfg.CNBOutboundEnabled,
    })

    if err == nil && len(mappings) > 0 && h.cfg.CNBOutboundEnabled {
        cn := &cnb.Client{BaseURL: h.cfg.CNBBaseURL, Token: h.cfg.CNBAppToken, IssueCreatePath: h.cfg.CNBIssueCreatePath, IssueUpdatePath: h.cfg.CNBIssueUpdatePath, IssueCommentPath: h.cfg.CNBIssueCommentPath}
        switch action {
        case "create", "created":
            // Fan-out create
            for _, m := range mappings {
                dir := strings.ToLower(m.SyncDirection.String)
                if !m.SyncDirection.Valid { dir = "" }
                hit := labelSelectorMatch(m.LabelSelector.String, labels)
                links, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID)
                exists := false
                for _, lk := range links { if lk.Repo == m.CNBRepoID { exists = true; break } }

                decision, skip := "create", ""
                if dir != "bidirectional" { decision, skip = "skip", "direction" }
                if decision != "skip" && !hit { decision, skip = "skip", "label_miss" }
                if decision != "skip" && exists { decision, skip = "skip", "already_linked" }

                LogStructured("info", map[string]any{"event": "plane.issue.route", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "direction": dir, "label_selector": m.LabelSelector.String, "label_hit": hit, "already_linked": exists, "decision": decision, "skip_reason": skip})
                if decision == "skip" { continue }

                unlock := h.acquireCreateLock(m.CNBRepoID, planeIssueID)
                existsLocked := false
                if links2, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID); len(links2) > 0 {
                    for _, lk := range links2 { if lk.Repo == m.CNBRepoID { existsLocked = true; break } }
                }
                if existsLocked {
                    unlock()
                    LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "op": "create_issue", "delivery_id": deliveryID, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "decision": "skip", "skip_reason": "already_linked"})
                    continue
                }
                iid, err := cn.CreateIssue(ctx, m.CNBRepoID, name, descHTML)
                if err != nil || iid == "" {
                    LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "op": "create_issue", "error": map[string]any{"code": "cnb_create_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                    unlock()
                    continue
                }
                _ = h.db.CreateIssueLink(ctx, planeIssueID, m.CNBRepoID, iid)
                LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "op": "create_issue", "result": "created", "cnb_issue_iid": iid})
                unlock()
                if startDate != "" || dueDate != "" || hasPriority {
                    dfields := map[string]any{}
                    if startDate != "" { dfields["start_date"] = startDate }
                    if dueDate != "" { dfields["end_date"] = dueDate }
                    if hasPriority { dfields["priority"] = cnbPriority }
                    LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "op": "update_issue", "fields_keys": keysOf(dfields)})
                    if err := cn.UpdateIssue(ctx, m.CNBRepoID, iid, dfields); err != nil {
                        LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "op": "update_issue", "error": map[string]any{"code": "cnb_update_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                    } else {
                        LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "op": "update_issue", "result": "updated"})
                    }
                }
                // Assignees
                if len(assigneePlaneIDs) > 0 {
                    if ids, _ := h.db.FindCNBUserIDsByPlaneUsers(ctx, assigneePlaneIDs); len(ids) > 0 {
                        LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "op": "update_assignees", "assignees_count": len(ids)})
                        if err := cn.UpdateIssueAssignees(ctx, m.CNBRepoID, iid, ids); err != nil {
                            LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "op": "update_assignees", "error": map[string]any{"code": "cnb_update_assignees_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                        } else {
                            LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "op": "update_assignees", "result": "updated"})
                        }
                    } else {
                        LogStructured("info", map[string]any{"event": "plane.issue.assignees.skip", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": m.CNBRepoID, "reason": "no_user_mapping", "plane_assignees": assigneePlaneIDs})
                    }
                }
                // Labels sync (best-effort)
                if len(labels) > 0 {
                    items := dataGetLabelItems(data)
                    toCreate := make([]cnb.Label, 0, len(items))
                    if len(items) > 0 {
                        for _, it := range items { toCreate = append(toCreate, cnb.Label{Name: it.Name, Color: strings.TrimPrefix(it.Color, "#")}) }
                    }
                    if len(toCreate) > 0 {
                        if err := cn.EnsureRepoLabelsWithColors(ctx, m.CNBRepoID, toCreate); err != nil {
                            LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "op": "ensure_labels", "repo": m.CNBRepoID, "delivery_id": deliveryID, "error": map[string]any{"code": "cnb_labels_ensure_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                        }
                    } else if err := cn.EnsureRepoLabels(ctx, m.CNBRepoID, labels); err != nil {
                        LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "op": "ensure_labels", "repo": m.CNBRepoID, "delivery_id": deliveryID, "error": map[string]any{"code": "cnb_labels_ensure_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                    } else {
                        if err := cn.SetIssueLabels(ctx, m.CNBRepoID, iid, labels); err != nil {
                            LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "op": "set_issue_labels", "repo": m.CNBRepoID, "delivery_id": deliveryID, "error": map[string]any{"code": "cnb_set_labels_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                        } else {
                            LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "op": "set_issue_labels", "repo": m.CNBRepoID, "delivery_id": deliveryID, "result": "set"})
                        }
                    }
                }
            }
        case "update", "updated":
            if links, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID); len(links) > 0 {
                fields := map[string]any{}
                if name != "" { fields["title"] = name }
                if descHTML != "" { fields["body"] = descHTML }
                if startDate != "" { fields["start_date"] = startDate }
                if dueDate != "" { fields["end_date"] = dueDate }
                // close if terminal state
                shouldClose := isTerminalPlaneState(stateName)
                if shouldClose {
                    fields["state"] = "closed"
                    fields["state_reason"] = cnbStateReasonForPlane(stateName)
                }
                // priority
                wantPriority := ""
                if hasPriority { fields["priority"] = cnbPriority; wantPriority = cnbPriority }
                var cnbAssignees []string
                if len(assigneePlaneIDs) > 0 {
                    if ids, _ := h.db.FindCNBUserIDsByPlaneUsers(ctx, assigneePlaneIDs); len(ids) > 0 { cnbAssignees = ids }
                }
                LogStructured("info", map[string]any{
                    "event":                 "plane.issue.update.decision",
                    "delivery_id":           deliveryID,
                    "action":                action,
                    "plane_issue_id":        planeIssueID,
                    "links_count":           len(links),
                    "title_set":             name != "",
                    "body_set":              descHTML != "",
                    "labels_in_payload":     len(labels),
                    "priority_set":          hasPriority,
                    "plane_priority":        planePriority,
                    "cnb_priority":          cnbPriority,
                    "start_date_set":        startDate != "",
                    "end_date_set":          dueDate != "",
                    "plane_state_name":      strings.ToLower(stateName),
                    "close_cnb":             shouldClose,
                    "plane_assignees_count": len(assigneePlaneIDs),
                    "cnb_assignees_count":   len(cnbAssignees),
                    "update_fields_count":   len(fields),
                })
                for _, lk := range links {
                    if len(fields) == 0 {
                        LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "update_issue", "decision": "skip", "skip_reason": "no_supported_fields"})
                        continue
                    }
                    keys := make([]string, 0, len(fields))
                    for k := range fields { keys = append(keys, k) }
                    LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "update_issue", "fields_keys": keys})
                    if err := cn.UpdateIssue(ctx, lk.Repo, lk.Number, fields); err != nil {
                        LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "update_issue", "error": map[string]any{"code": "cnb_update_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                    } else {
                        LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "update_issue", "result": "updated"})
                        // verify priority if needed, with fallback single-field patch
                        if wantPriority != "" {
                            if det, err := cn.GetIssue(ctx, lk.Repo, lk.Number); err == nil {
                                if !strings.EqualFold(strings.TrimSpace(det.Priority), strings.TrimSpace(wantPriority)) {
                                    LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "verify_priority", "error": map[string]any{"code": "cnb_priority_not_applied", "message": "priority mismatch after patch", "want": wantPriority, "got": det.Priority}})
                                    if err := cn.UpdateIssue(ctx, lk.Repo, lk.Number, map[string]any{"priority": wantPriority}); err != nil {
                                        LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "update_issue", "error": map[string]any{"code": "cnb_priority_patch_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                                    } else if det2, err2 := cn.GetIssue(ctx, lk.Repo, lk.Number); err2 == nil {
                                        LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "verify_priority", "result": map[string]any{"want": wantPriority, "got": det2.Priority}})
                                    }
                                } else {
                                    LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "verify_priority", "result": "ok", "priority": det.Priority})
                                }
                            }
                        }
                    }
                    // Assignees after update
                    if len(cnbAssignees) > 0 {
                        LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "update_assignees", "assignees_count": len(cnbAssignees)})
                        if err := cn.UpdateIssueAssignees(ctx, lk.Repo, lk.Number, cnbAssignees); err != nil {
                            LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "update_assignees", "error": map[string]any{"code": "cnb_update_assignees_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                        } else {
                            LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "update_assignees", "result": "updated"})
                        }
                    } else if len(assigneePlaneIDs) > 0 {
                        LogStructured("info", map[string]any{"event": "plane.issue.assignees.skip", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "reason": "no_user_mapping", "plane_assignees": assigneePlaneIDs})
                    }
                }
                // labels best-effort
                if len(labels) > 0 {
                    items := dataGetLabelItems(data)
                    toCreate := make([]cnb.Label, 0, len(items))
                    if len(items) > 0 {
                        for _, it := range items { toCreate = append(toCreate, cnb.Label{Name: it.Name, Color: strings.TrimPrefix(it.Color, "#")}) }
                    }
                    if len(toCreate) > 0 {
                        if err := cn.EnsureRepoLabelsWithColors(ctx, links[0].Repo, toCreate); err != nil {
                            LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "op": "ensure_labels", "repo": links[0].Repo, "delivery_id": deliveryID, "error": map[string]any{"code": "cnb_labels_ensure_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                        }
                    } else if err := cn.EnsureRepoLabels(ctx, links[0].Repo, labels); err != nil {
                        LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "op": "ensure_labels", "repo": links[0].Repo, "delivery_id": deliveryID, "error": map[string]any{"code": "cnb_labels_ensure_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                    } else {
                        for _, lk := range links {
                            if err := cn.SetIssueLabels(ctx, lk.Repo, lk.Number, labels); err != nil {
                                LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "op": "set_issue_labels", "repo": lk.Repo, "delivery_id": deliveryID, "error": map[string]any{"code": "cnb_set_labels_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                            } else {
                                LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "op": "set_issue_labels", "repo": lk.Repo, "delivery_id": deliveryID, "result": "set"})
                            }
                        }
                    }
                }
            }
        case "delete", "deleted", "close", "closed":
            if links, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID); len(links) > 0 {
                for _, lk := range links {
                    if err := cn.CloseIssue(ctx, lk.Repo, lk.Number); err != nil {
                        LogStructured("error", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "close_issue", "error": map[string]any{"code": "cnb_close_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                    } else {
                        LogStructured("info", map[string]any{"event": "plane.issue.cnbrpc", "delivery_id": deliveryID, "action": action, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "close_issue", "result": "closed"})
                    }
                }
            }
        }
    }

    // Notify Feishu thread if bound, but avoid duplicate notifications when this issue event is caused by a comment activity.
    // Plane 会同时发送 issue_comment 与 issue.updated(仅评论相关字段变更) 两个事件。为避免重复通知，这里在 activity.field 含 "comment" 时不发送摘要通知。
    af := strings.ToLower(strings.TrimSpace(env.Activity.Field))
    isCommentActivity := strings.Contains(af, "comment")
    if !isCommentActivity && planeIssueID != "" {
        if tid, err := h.db.FindLarkThreadByPlaneIssue(ctx, planeIssueID); err == nil && tid != "" {
            summary := "Plane 工作项更新: " + truncate(name, 80)
            if action != "" { summary += " (" + action + ")" }
            go h.sendLarkTextToThread("", tid, summary)
        }
    }
    if deliveryID != "" {
        _ = h.db.UpdateEventDeliveryStatus(ctx, "plane.issue", deliveryID, "succeeded", nil)
    }
}

func (h *Handler) handlePlaneIssueComment(env planeWebhookEnvelope, deliveryID string) {
    if !hHasDB(h) { return }
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    data := env.Data
    planeIssueID, _ := dataGetString(data, "issue")
    commentHTML, _ := dataGetString(data, "comment_html")
    if planeIssueID == "" || commentHTML == "" { return }
    if h.cfg.CNBOutboundEnabled {
        if links, _ := h.db.ListCNBIssuesByPlaneIssue(ctx, planeIssueID); len(links) > 0 {
            cn := &cnb.Client{BaseURL: h.cfg.CNBBaseURL, Token: h.cfg.CNBAppToken, IssueCreatePath: h.cfg.CNBIssueCreatePath, IssueUpdatePath: h.cfg.CNBIssueUpdatePath, IssueCommentPath: h.cfg.CNBIssueCommentPath}
            for _, lk := range links {
                if err := cn.AddComment(ctx, lk.Repo, lk.Number, commentHTML); err != nil {
                    LogStructured("error", map[string]any{"event": "plane.issue_comment.cnbrpc", "delivery_id": deliveryID, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "add_comment", "error": map[string]any{"code": "cnb_comment_failed", "message": truncate(fmt.Sprintf("%v", err), 200)}})
                } else {
                    LogStructured("info", map[string]any{"event": "plane.issue_comment.cnbrpc", "delivery_id": deliveryID, "plane_issue_id": planeIssueID, "repo": lk.Repo, "op": "add_comment", "result": "commented"})
                }
            }
        }
    }
    if tid, err := h.db.FindLarkThreadByPlaneIssue(ctx, planeIssueID); err == nil && tid != "" {
        txt := commentHTML
        txt = strings.ReplaceAll(txt, "<br>", "\n")
        txt = stripTags(txt)
        msg := "Plane 评论: " + truncate(txt, 140)
        go h.sendLarkTextToThread("", tid, msg)
    }
    if deliveryID != "" {
        _ = h.db.UpdateEventDeliveryStatus(ctx, "plane.issue_comment", deliveryID, "succeeded", nil)
    }
}

// ==== helpers (issue) ====

func hHasDB(h *Handler) bool { return h != nil && h.db != nil && h.db.SQL != nil }

func dataGetString(m map[string]any, key string) (string, bool) {
    if m == nil { return "", false }
    if v, ok := m[key]; ok {
        if s, ok := v.(string); ok { return s, true }
    }
    return "", false
}

func dataGetLabels(m map[string]any) []string {
    names := make([]string, 0, 8)
    if m == nil { return names }
    if v, ok := m["labels"]; ok {
        if arr, ok := v.([]any); ok {
            for _, it := range arr {
                switch t := it.(type) {
                case map[string]any:
                    if n, ok := t["name"].(string); ok && n != "" { names = append(names, n) }
                case string:
                    if t != "" { names = append(names, t) }
                }
            }
        }
    }
    if len(names) == 0 {
        if v, ok := m["label_names"]; ok {
            if arr, ok := v.([]any); ok {
                for _, it := range arr { if s, ok := it.(string); ok && s != "" { names = append(names, s) } }
            }
        }
    }
    return names
}

func dataGetDates(m map[string]any) (start string, due string) {
    if m == nil { return "", "" }
    readFirst := func(keys ...string) string {
        for _, k := range keys {
            if v, ok := m[k]; ok {
                if s, ok := v.(string); ok && strings.TrimSpace(s) != "" { return s }
            }
        }
        return ""
    }
    start = readFirst("start_date", "start_on", "start", "start_at")
    due = readFirst("due_date", "target_date", "due_on", "due", "target_on", "target_at")
    return
}

// Priority mapping: urgent/high/medium/low/none or numeric 4..0 to CNB string P0..P3/""
func dataGetPriority(m map[string]any) (string, string, bool) {
    if m == nil { return "", "", false }
    tryKeys := []string{"priority", "priority_name", "priority_label", "priority_value", "priority_level"}
    var raw any
    for _, k := range tryKeys { if v, ok := m[k]; ok { raw = v; break } }
    if raw == nil { return "", "", false }
    strMap := func(s string) (string, string, bool) {
        n := strings.ToLower(strings.TrimSpace(s))
        if n == "" { return "", "", false }
        switch n {
        case "urgent", "critical", "blocker", "p0":
            return "P0", "urgent", true
        case "high", "p1":
            return "P1", "high", true
        case "medium", "normal", "p2":
            return "P2", "medium", true
        case "low", "minor", "p3":
            return "P3", "low", true
        case "none", "no", "null", "none_priority", "":
            return "", "none", true
        case "-1p":
            return "-1P", "-1p", true
        case "-2p":
            return "-2P", "-2p", true
        }
        up := strings.ToUpper(n)
        if up == "P0" || up == "P1" || up == "P2" || up == "P3" || up == "-1P" || up == "-2P" { return up, n, true }
        return "", n, false
    }
    numMap := func(f float64) (string, string, bool) {
        switch int(f) {
        case 4:
            return "P0", "urgent", true
        case 3:
            return "P1", "high", true
        case 2:
            return "P2", "medium", true
        case 1:
            return "P3", "low", true
        case 0:
            return "", "none", true
        }
        return "", "", false
    }
    switch t := raw.(type) {
    case string:
        return strMap(t)
    case float64:
        return numMap(t)
    case int:
        return numMap(float64(t))
    case map[string]any:
        if v, ok := t["name"]; ok { if s, ok := v.(string); ok { return strMap(s) } }
        if v, ok := t["value"]; ok {
            switch vv := v.(type) {
            case string:
                return strMap(vv)
            case float64:
                return numMap(vv)
            }
        }
    }
    return "", "", false
}

func keysOf(m map[string]any) []string {
    keys := make([]string, 0, len(m))
    for k := range m { keys = append(keys, k) }
    sort.Strings(keys)
    return keys
}

// dataGetAssigneeIDs extracts Plane user IDs from webhook payload.
func dataGetAssigneeIDs(m map[string]any) []string {
    ids := make([]string, 0, 4)
    if m == nil { return ids }
    add := func(s string) {
        s = strings.TrimSpace(s)
        if s == "" { return }
        for _, e := range ids { if e == s { return } }
        ids = append(ids, s)
    }
    tryAppend := func(v any) {
        switch t := v.(type) {
        case []any:
            for _, it := range t {
                switch a := it.(type) {
                case string:
                    add(a)
                case map[string]any:
                    if s, ok := a["id"].(string); ok { add(s) }
                    if s, ok := a["user_id"].(string); ok { add(s) }
                }
            }
        case []string:
            for _, s := range t { add(s) }
        }
    }
    if v, ok := m["assignees"]; ok { tryAppend(v) }
    if v, ok := m["assignee_ids"]; ok { tryAppend(v) }
    return ids
}

type planeLabel struct{ Name, Color string }

func dataGetLabelItems(m map[string]any) []planeLabel {
    out := make([]planeLabel, 0, 8)
    if m == nil { return out }
    if v, ok := m["labels"]; ok {
        if arr, ok := v.([]any); ok {
            for _, it := range arr {
                if mp, ok := it.(map[string]any); ok {
                    name, _ := mp["name"].(string)
                    color, _ := mp["color"].(string)
                    if strings.TrimSpace(name) != "" { out = append(out, planeLabel{Name: name, Color: color}) }
                }
            }
        }
    }
    return out
}

func labelSelectorMatch(selector string, labels []string) bool {
    selector = strings.TrimSpace(selector)
    if selector == "" { return false }
    tokens := make([]string, 0, 8)
    for _, p := range strings.FieldsFunc(selector, func(r rune) bool { return r == ',' || r == ' ' || r == ';' || r == '|' }) {
        p = strings.TrimSpace(p)
        if p != "" { tokens = append(tokens, strings.ToLower(p)) }
    }
    if len(tokens) == 0 { return false }
    for _, t := range tokens { if t == "*" || t == "all" { return true } }
    set := make(map[string]struct{}, len(labels))
    for _, l := range labels { if l != "" { set[strings.ToLower(l)] = struct{}{} } }
    for _, tok := range tokens { if _, ok := set[tok]; ok { return true } }
    return false
}

func dataGetStateNameAndID(m map[string]any) (name string, id string) {
    if m == nil { return "", "" }
    if v, ok := m["state"]; ok {
        switch t := v.(type) {
        case map[string]any:
            if s, ok := t["name"].(string); ok { name = s }
            if s, ok := t["id"].(string); ok { id = s }
        case string:
            id = t
        }
    }
    return strings.TrimSpace(name), strings.TrimSpace(id)
}

func anyToStateName(v any) (string, bool) {
    switch t := v.(type) {
    case string:
        s := strings.TrimSpace(t)
        if s == "" { return "", false }
        return s, true
    case map[string]any:
        if s, ok := t["name"].(string); ok && strings.TrimSpace(s) != "" { return s, true }
    }
    return "", false
}

func isTerminalPlaneState(name string) bool {
    n := strings.ToLower(strings.TrimSpace(name))
    if n == "" { return false }
    switch n {
    case "done", "completed", "complete", "cancelled", "canceled":
        return true
    }
    return false
}

// cnbStateReasonForPlane maps Plane terminal state name to CNB state_reason.
// CNB allowed: completed, not_planned, reopened.
func cnbStateReasonForPlane(name string) string {
    n := strings.ToLower(strings.TrimSpace(name))
    if strings.Contains(n, "cancel") {
        return "not_planned"
    }
    return "completed"
}
