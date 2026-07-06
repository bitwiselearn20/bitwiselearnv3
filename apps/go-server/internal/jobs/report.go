package jobs

// AssessmentReportQueue is the durable queue name for triggered assessment
// report generation (ports the publish_message("assessment-report", ...)
// call in routers/assessment.py). Consumed by apps/python-worker today;
// the Go worker doesn't process this queue yet — Phase 3 only wires the
// publish side so the trigger endpoint keeps its existing side effect.
const AssessmentReportQueue = "assessment-report"

// AssessmentReportJob is the message body for AssessmentReportQueue.
type AssessmentReportJob struct {
	AssessmentID string `json:"assessment_id"`
}
