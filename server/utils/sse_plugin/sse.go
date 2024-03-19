package sse_plugin

import "github.com/gin-gonic/gin"

type SSEHandler func(ctx *gin.Context)
