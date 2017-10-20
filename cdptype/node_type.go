package cdptype

import (
	"github.com/mafredri/cdp/protocol"
)

const (
	NodeTypeElement               int = 1
	NodeTypeAttribute                 = 2
	NodeTypeText                      = 3
	NodeTypeCDATA                     = 4
	NodeTypeEntityReference           = 5
	NodeTypeEntity                    = 6
	NodeTypeProcessingInstruction     = 7
	NodeTypeComment                   = 8
	NodeTypeDocument                  = 9
	NodeTypeDocumentType              = 10
	NodeTypeDocumentFragment          = 11
	NodeTypeNotation                  = 12
)

const (
	ResourceTypeDocument    protocol.PageResourceType = "Document"
	ResourceTypeStylesheet                            = "Stylesheet"
	ResourceTypeImage                                 = "Image"
	ResourceTypeMedia                                 = "Media"
	ResourceTypeFont                                  = "Font"
	ResourceTypeScript                                = "Script"
	ResourceTypeTextTrack                             = "TextTrack"
	ResourceTypeXHR                                   = "XHR"
	ResourceTypeFetch                                 = "Fetch"
	ResourceTypeEventSource                           = "EventSource"
	ResourceTypeWebSocket                             = "WebSocket"
	ResourceTypeManifest                              = "Manifest"
	ResourceTypeOther                                 = "Other"
)
