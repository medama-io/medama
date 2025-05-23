// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ogen-go/ogen/uri"
)

func (s *Server) cutPrefix(path string) (string, bool) {
	prefix := s.cfg.Prefix
	if prefix == "" {
		return path, true
	}
	if !strings.HasPrefix(path, prefix) {
		// Prefix doesn't match.
		return "", false
	}
	// Cut prefix from the path.
	return strings.TrimPrefix(path, prefix), true
}

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Path
	elemIsEscaped := false
	if rawPath := r.URL.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
			elemIsEscaped = strings.ContainsRune(elem, '%')
		}
	}

	elem, ok := s.cutPrefix(elem)
	if !ok || len(elem) == 0 {
		s.notFound(w, r)
		return
	}
	args := [1]string{}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"

			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'a': // Prefix: "auth/log"

				if l := len("auth/log"); len(elem) >= l && elem[0:l] == "auth/log" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'i': // Prefix: "in"

					if l := len("in"); len(elem) >= l && elem[0:l] == "in" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handlePostAuthLoginRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

				case 'o': // Prefix: "out"

					if l := len("out"); len(elem) >= l && elem[0:l] == "out" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handlePostAuthLogoutRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

				}

			case 'e': // Prefix: "event/"

				if l := len("event/"); len(elem) >= l && elem[0:l] == "event/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'h': // Prefix: "hit"

					if l := len("hit"); len(elem) >= l && elem[0:l] == "hit" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handlePostEventHitRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

				case 'p': // Prefix: "ping"

					if l := len("ping"); len(elem) >= l && elem[0:l] == "ping" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleGetEventPingRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}

				}

			case 'u': // Prefix: "user"

				if l := len("user"); len(elem) >= l && elem[0:l] == "user" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch r.Method {
					case "DELETE":
						s.handleDeleteUserRequest([0]string{}, elemIsEscaped, w, r)
					case "GET":
						s.handleGetUserRequest([0]string{}, elemIsEscaped, w, r)
					case "PATCH":
						s.handlePatchUserRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "DELETE,GET,PATCH")
					}

					return
				}
				switch elem[0] {
				case '/': // Prefix: "/usage"

					if l := len("/usage"); len(elem) >= l && elem[0:l] == "/usage" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleGetUserUsageRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}

				}

			case 'w': // Prefix: "website"

				if l := len("website"); len(elem) >= l && elem[0:l] == "website" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/"

					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "hostname"
					// Match until "/"
					idx := strings.IndexByte(elem, '/')
					if idx < 0 {
						idx = len(elem)
					}
					args[0] = elem[:idx]
					elem = elem[idx:]

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"

						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							break
						}
						switch elem[0] {
						case 'b': // Prefix: "browsers"

							if l := len("browsers"); len(elem) >= l && elem[0:l] == "browsers" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleGetWebsiteIDBrowsersRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}

						case 'c': // Prefix: "c"

							if l := len("c"); len(elem) >= l && elem[0:l] == "c" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'a': // Prefix: "ampaigns"

								if l := len("ampaigns"); len(elem) >= l && elem[0:l] == "ampaigns" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetWebsiteIDCampaignsRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}

							case 'o': // Prefix: "ountries"

								if l := len("ountries"); len(elem) >= l && elem[0:l] == "ountries" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetWebsiteIDCountryRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}

							}

						case 'd': // Prefix: "devices"

							if l := len("devices"); len(elem) >= l && elem[0:l] == "devices" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleGetWebsiteIDDeviceRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}

						case 'l': // Prefix: "languages"

							if l := len("languages"); len(elem) >= l && elem[0:l] == "languages" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleGetWebsiteIDLanguageRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}

						case 'm': // Prefix: "mediums"

							if l := len("mediums"); len(elem) >= l && elem[0:l] == "mediums" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleGetWebsiteIDMediumsRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}

						case 'o': // Prefix: "os"

							if l := len("os"); len(elem) >= l && elem[0:l] == "os" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleGetWebsiteIDOsRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}

						case 'p': // Prefix: "p"

							if l := len("p"); len(elem) >= l && elem[0:l] == "p" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'a': // Prefix: "ages"

								if l := len("ages"); len(elem) >= l && elem[0:l] == "ages" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetWebsiteIDPagesRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}

							case 'r': // Prefix: "roperties"

								if l := len("roperties"); len(elem) >= l && elem[0:l] == "roperties" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetWebsiteIDPropertiesRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}

							}

						case 'r': // Prefix: "referrers"

							if l := len("referrers"); len(elem) >= l && elem[0:l] == "referrers" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleGetWebsiteIDReferrersRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}

						case 's': // Prefix: "s"

							if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'o': // Prefix: "ources"

								if l := len("ources"); len(elem) >= l && elem[0:l] == "ources" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetWebsiteIDSourcesRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}

							case 'u': // Prefix: "ummary"

								if l := len("ummary"); len(elem) >= l && elem[0:l] == "ummary" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetWebsiteIDSummaryRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}

							}

						case 't': // Prefix: "time"

							if l := len("time"); len(elem) >= l && elem[0:l] == "time" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleGetWebsiteIDTimeRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}

						}

					}

				case 's': // Prefix: "s"

					if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch r.Method {
						case "GET":
							s.handleGetWebsitesRequest([0]string{}, elemIsEscaped, w, r)
						case "POST":
							s.handlePostWebsitesRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET,POST")
						}

						return
					}
					switch elem[0] {
					case '/': // Prefix: "/"

						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "hostname"
						// Leaf parameter, slashes are prohibited
						idx := strings.IndexByte(elem, '/')
						if idx >= 0 {
							break
						}
						args[0] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "DELETE":
								s.handleDeleteWebsitesIDRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							case "GET":
								s.handleGetWebsitesIDRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							case "PATCH":
								s.handlePatchWebsitesIDRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "DELETE,GET,PATCH")
							}

							return
						}

					}

				}

			}

		}
	}
	s.notFound(w, r)
}

// Route is route object.
type Route struct {
	name        string
	summary     string
	operationID string
	pathPattern string
	count       int
	args        [1]string
}

// Name returns ogen operation name.
//
// It is guaranteed to be unique and not empty.
func (r Route) Name() string {
	return r.name
}

// Summary returns OpenAPI summary.
func (r Route) Summary() string {
	return r.summary
}

// OperationID returns OpenAPI operationId.
func (r Route) OperationID() string {
	return r.operationID
}

// PathPattern returns OpenAPI path.
func (r Route) PathPattern() string {
	return r.pathPattern
}

// Args returns parsed arguments.
func (r Route) Args() []string {
	return r.args[:r.count]
}

// FindRoute finds Route for given method and path.
//
// Note: this method does not unescape path or handle reserved characters in path properly. Use FindPath instead.
func (s *Server) FindRoute(method, path string) (Route, bool) {
	return s.FindPath(method, &url.URL{Path: path})
}

// FindPath finds Route for given method and URL.
func (s *Server) FindPath(method string, u *url.URL) (r Route, _ bool) {
	var (
		elem = u.Path
		args = r.args
	)
	if rawPath := u.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
		}
		defer func() {
			for i, arg := range r.args[:r.count] {
				if unescaped, err := url.PathUnescape(arg); err == nil {
					r.args[i] = unescaped
				}
			}
		}()
	}

	elem, ok := s.cutPrefix(elem)
	if !ok {
		return r, false
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"

			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'a': // Prefix: "auth/log"

				if l := len("auth/log"); len(elem) >= l && elem[0:l] == "auth/log" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'i': // Prefix: "in"

					if l := len("in"); len(elem) >= l && elem[0:l] == "in" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "POST":
							r.name = PostAuthLoginOperation
							r.summary = "Login"
							r.operationID = "post-auth-login"
							r.pathPattern = "/auth/login"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

				case 'o': // Prefix: "out"

					if l := len("out"); len(elem) >= l && elem[0:l] == "out" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "POST":
							r.name = PostAuthLogoutOperation
							r.summary = "Logout"
							r.operationID = "post-auth-logout"
							r.pathPattern = "/auth/logout"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

				}

			case 'e': // Prefix: "event/"

				if l := len("event/"); len(elem) >= l && elem[0:l] == "event/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'h': // Prefix: "hit"

					if l := len("hit"); len(elem) >= l && elem[0:l] == "hit" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "POST":
							r.name = PostEventHitOperation
							r.summary = "Send Hit Event"
							r.operationID = "post-event-hit"
							r.pathPattern = "/event/hit"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

				case 'p': // Prefix: "ping"

					if l := len("ping"); len(elem) >= l && elem[0:l] == "ping" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "GET":
							r.name = GetEventPingOperation
							r.summary = "Ping"
							r.operationID = "get-event-ping"
							r.pathPattern = "/event/ping"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

				}

			case 'u': // Prefix: "user"

				if l := len("user"); len(elem) >= l && elem[0:l] == "user" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch method {
					case "DELETE":
						r.name = DeleteUserOperation
						r.summary = "Delete User"
						r.operationID = "delete-user"
						r.pathPattern = "/user"
						r.args = args
						r.count = 0
						return r, true
					case "GET":
						r.name = GetUserOperation
						r.summary = "Get User Info"
						r.operationID = "get-user"
						r.pathPattern = "/user"
						r.args = args
						r.count = 0
						return r, true
					case "PATCH":
						r.name = PatchUserOperation
						r.summary = "Update User Info"
						r.operationID = "patch-user"
						r.pathPattern = "/user"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}
				switch elem[0] {
				case '/': // Prefix: "/usage"

					if l := len("/usage"); len(elem) >= l && elem[0:l] == "/usage" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "GET":
							r.name = GetUserUsageOperation
							r.summary = "Get Resource Usage"
							r.operationID = "get-user-usage"
							r.pathPattern = "/user/usage"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

				}

			case 'w': // Prefix: "website"

				if l := len("website"); len(elem) >= l && elem[0:l] == "website" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/"

					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "hostname"
					// Match until "/"
					idx := strings.IndexByte(elem, '/')
					if idx < 0 {
						idx = len(elem)
					}
					args[0] = elem[:idx]
					elem = elem[idx:]

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"

						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							break
						}
						switch elem[0] {
						case 'b': // Prefix: "browsers"

							if l := len("browsers"); len(elem) >= l && elem[0:l] == "browsers" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "GET":
									r.name = GetWebsiteIDBrowsersOperation
									r.summary = "Get Browser Stats"
									r.operationID = "get-website-id-browsers"
									r.pathPattern = "/website/{hostname}/browsers"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

						case 'c': // Prefix: "c"

							if l := len("c"); len(elem) >= l && elem[0:l] == "c" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'a': // Prefix: "ampaigns"

								if l := len("ampaigns"); len(elem) >= l && elem[0:l] == "ampaigns" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "GET":
										r.name = GetWebsiteIDCampaignsOperation
										r.summary = "Get UTM Campaign Stats"
										r.operationID = "get-website-id-campaigns"
										r.pathPattern = "/website/{hostname}/campaigns"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}

							case 'o': // Prefix: "ountries"

								if l := len("ountries"); len(elem) >= l && elem[0:l] == "ountries" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "GET":
										r.name = GetWebsiteIDCountryOperation
										r.summary = "Get Country Stats"
										r.operationID = "get-website-id-country"
										r.pathPattern = "/website/{hostname}/countries"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}

							}

						case 'd': // Prefix: "devices"

							if l := len("devices"); len(elem) >= l && elem[0:l] == "devices" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "GET":
									r.name = GetWebsiteIDDeviceOperation
									r.summary = "Get Device Stats"
									r.operationID = "get-website-id-device"
									r.pathPattern = "/website/{hostname}/devices"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

						case 'l': // Prefix: "languages"

							if l := len("languages"); len(elem) >= l && elem[0:l] == "languages" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "GET":
									r.name = GetWebsiteIDLanguageOperation
									r.summary = "Get Language Stats"
									r.operationID = "get-website-id-language"
									r.pathPattern = "/website/{hostname}/languages"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

						case 'm': // Prefix: "mediums"

							if l := len("mediums"); len(elem) >= l && elem[0:l] == "mediums" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "GET":
									r.name = GetWebsiteIDMediumsOperation
									r.summary = "Get UTM Medium Stats"
									r.operationID = "get-website-id-mediums"
									r.pathPattern = "/website/{hostname}/mediums"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

						case 'o': // Prefix: "os"

							if l := len("os"); len(elem) >= l && elem[0:l] == "os" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "GET":
									r.name = GetWebsiteIDOsOperation
									r.summary = "Get OS Stats"
									r.operationID = "get-website-id-os"
									r.pathPattern = "/website/{hostname}/os"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

						case 'p': // Prefix: "p"

							if l := len("p"); len(elem) >= l && elem[0:l] == "p" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'a': // Prefix: "ages"

								if l := len("ages"); len(elem) >= l && elem[0:l] == "ages" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "GET":
										r.name = GetWebsiteIDPagesOperation
										r.summary = "Get Page Stats"
										r.operationID = "get-website-id-pages"
										r.pathPattern = "/website/{hostname}/pages"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}

							case 'r': // Prefix: "roperties"

								if l := len("roperties"); len(elem) >= l && elem[0:l] == "roperties" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "GET":
										r.name = GetWebsiteIDPropertiesOperation
										r.summary = "Get Property Stats"
										r.operationID = "get-website-id-properties"
										r.pathPattern = "/website/{hostname}/properties"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}

							}

						case 'r': // Prefix: "referrers"

							if l := len("referrers"); len(elem) >= l && elem[0:l] == "referrers" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "GET":
									r.name = GetWebsiteIDReferrersOperation
									r.summary = "Get Referrer Stats"
									r.operationID = "get-website-id-referrers"
									r.pathPattern = "/website/{hostname}/referrers"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

						case 's': // Prefix: "s"

							if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'o': // Prefix: "ources"

								if l := len("ources"); len(elem) >= l && elem[0:l] == "ources" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "GET":
										r.name = GetWebsiteIDSourcesOperation
										r.summary = "Get UTM Source Stats"
										r.operationID = "get-website-id-sources"
										r.pathPattern = "/website/{hostname}/sources"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}

							case 'u': // Prefix: "ummary"

								if l := len("ummary"); len(elem) >= l && elem[0:l] == "ummary" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "GET":
										r.name = GetWebsiteIDSummaryOperation
										r.summary = "Get Stat Summary"
										r.operationID = "get-website-id-summary"
										r.pathPattern = "/website/{hostname}/summary"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}

							}

						case 't': // Prefix: "time"

							if l := len("time"); len(elem) >= l && elem[0:l] == "time" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "GET":
									r.name = GetWebsiteIDTimeOperation
									r.summary = "Get Time Stats"
									r.operationID = "get-website-id-time"
									r.pathPattern = "/website/{hostname}/time"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

						}

					}

				case 's': // Prefix: "s"

					if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							r.name = GetWebsitesOperation
							r.summary = "List Websites"
							r.operationID = "get-websites"
							r.pathPattern = "/websites"
							r.args = args
							r.count = 0
							return r, true
						case "POST":
							r.name = PostWebsitesOperation
							r.summary = "Add Website"
							r.operationID = "post-websites"
							r.pathPattern = "/websites"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
					switch elem[0] {
					case '/': // Prefix: "/"

						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "hostname"
						// Leaf parameter, slashes are prohibited
						idx := strings.IndexByte(elem, '/')
						if idx >= 0 {
							break
						}
						args[0] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "DELETE":
								r.name = DeleteWebsitesIDOperation
								r.summary = "Delete Website"
								r.operationID = "delete-websites-id"
								r.pathPattern = "/websites/{hostname}"
								r.args = args
								r.count = 1
								return r, true
							case "GET":
								r.name = GetWebsitesIDOperation
								r.summary = "Get Website"
								r.operationID = "get-websites-id"
								r.pathPattern = "/websites/{hostname}"
								r.args = args
								r.count = 1
								return r, true
							case "PATCH":
								r.name = PatchWebsitesIDOperation
								r.summary = "Update Website"
								r.operationID = "patch-websites-id"
								r.pathPattern = "/websites/{hostname}"
								r.args = args
								r.count = 1
								return r, true
							default:
								return
							}
						}

					}

				}

			}

		}
	}
	return r, false
}
