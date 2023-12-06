package model

import "github.com/yudeguang/ratelimit"

type RateLimitRule struct {
	IPRole    *ratelimit.Rule
	VisitRole *ratelimit.Rule
}
