package logger

// MarkerSilent is the logrus field key set by Logger.SilentError to mark an
// entry as "logged but should not trigger external alert sinks". Alert hooks
// (e.g. Sentry) MUST check this key and skip the entry when its value is true.
//
// The key is intentionally prefixed with an underscore so it is visually
// distinct from real payload fields and unlikely to collide with caller keys.
const MarkerSilent = "_silent"
