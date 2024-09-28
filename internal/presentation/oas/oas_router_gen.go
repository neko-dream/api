// Code generated by ogen, DO NOT EDIT.

package oas

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
	args := [2]string{}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/a"
			origElem := elem
			if l := len("/a"); len(elem) >= l && elem[0:l] == "/a" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'p': // Prefix: "pi/"
				origElem := elem
				if l := len("pi/"); len(elem) >= l && elem[0:l] == "pi/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 't': // Prefix: "talksession"
					origElem := elem
					if l := len("talksession"); len(elem) >= l && elem[0:l] == "talksession" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						origElem := elem
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "talkSessionID"
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
						case '/': // Prefix: "/opinions"
							origElem := elem
							if l := len("/opinions"); len(elem) >= l && elem[0:l] == "/opinions" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleListOpinionsRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}

							elem = origElem
						}

						elem = origElem
					case 's': // Prefix: "s"
						origElem := elem
						if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch r.Method {
							case "GET":
								s.handleGetTalkSessionsRequest([0]string{}, elemIsEscaped, w, r)
							case "POST":
								s.handleCreateTalkSessionRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET,POST")
							}

							return
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							origElem := elem
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "talkSessionId"
							// Match until "/"
							idx := strings.IndexByte(elem, '/')
							if idx < 0 {
								idx = len(elem)
							}
							args[0] = elem[:idx]
							elem = elem[idx:]

							if len(elem) == 0 {
								switch r.Method {
								case "GET":
									s.handleGetTalkSessionDetailRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}
							switch elem[0] {
							case '/': // Prefix: "/opinions"
								origElem := elem
								if l := len("/opinions"); len(elem) >= l && elem[0:l] == "/opinions" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									switch r.Method {
									case "POST":
										s.handlePostOpinionPostRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "POST")
									}

									return
								}
								switch elem[0] {
								case '/': // Prefix: "/"
									origElem := elem
									if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
										elem = elem[l:]
									} else {
										break
									}

									// Param: "opinionID"
									// Match until "/"
									idx := strings.IndexByte(elem, '/')
									if idx < 0 {
										idx = len(elem)
									}
									args[1] = elem[:idx]
									elem = elem[idx:]

									if len(elem) == 0 {
										break
									}
									switch elem[0] {
									case '/': // Prefix: "/intentions"
										origElem := elem
										if l := len("/intentions"); len(elem) >= l && elem[0:l] == "/intentions" {
											elem = elem[l:]
										} else {
											break
										}

										if len(elem) == 0 {
											// Leaf node.
											switch r.Method {
											case "POST":
												s.handleIntentionRequest([2]string{
													args[0],
													args[1],
												}, elemIsEscaped, w, r)
											default:
												s.notAllowed(w, r, "POST")
											}

											return
										}

										elem = origElem
									}

									elem = origElem
								}

								elem = origElem
							}

							elem = origElem
						}

						elem = origElem
					}

					elem = origElem
				case 'u': // Prefix: "user"
					origElem := elem
					if l := len("user"); len(elem) >= l && elem[0:l] == "user" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch r.Method {
						case "GET":
							s.handleGetUserProfileRequest([0]string{}, elemIsEscaped, w, r)
						case "PUT":
							s.handleEditUserProfileRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET,PUT")
						}

						return
					}
					switch elem[0] {
					case '/': // Prefix: "/register"
						origElem := elem
						if l := len("/register"); len(elem) >= l && elem[0:l] == "/register" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleRegisterUserRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'u': // Prefix: "uth/"
				origElem := elem
				if l := len("uth/"); len(elem) >= l && elem[0:l] == "uth/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "provider"
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
					origElem := elem
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'c': // Prefix: "callback"
						origElem := elem
						if l := len("callback"); len(elem) >= l && elem[0:l] == "callback" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleOAuthCallbackRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}

						elem = origElem
					case 'l': // Prefix: "login"
						origElem := elem
						if l := len("login"); len(elem) >= l && elem[0:l] == "login" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleAuthLoginRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			}

			elem = origElem
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
	args        [2]string
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
		case '/': // Prefix: "/a"
			origElem := elem
			if l := len("/a"); len(elem) >= l && elem[0:l] == "/a" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'p': // Prefix: "pi/"
				origElem := elem
				if l := len("pi/"); len(elem) >= l && elem[0:l] == "pi/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 't': // Prefix: "talksession"
					origElem := elem
					if l := len("talksession"); len(elem) >= l && elem[0:l] == "talksession" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						origElem := elem
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "talkSessionID"
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
						case '/': // Prefix: "/opinions"
							origElem := elem
							if l := len("/opinions"); len(elem) >= l && elem[0:l] == "/opinions" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								// Leaf node.
								switch method {
								case "GET":
									r.name = "ListOpinions"
									r.summary = "セッションの意見一覧"
									r.operationID = "listOpinions"
									r.pathPattern = "/api/talksession/{talkSessionID}/opinions"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}

							elem = origElem
						}

						elem = origElem
					case 's': // Prefix: "s"
						origElem := elem
						if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "GET":
								r.name = "GetTalkSessions"
								r.summary = "トークセッションリスト"
								r.operationID = "getTalkSessions"
								r.pathPattern = "/api/talksessions"
								r.args = args
								r.count = 0
								return r, true
							case "POST":
								r.name = "CreateTalkSession"
								r.summary = "トークセッション作成"
								r.operationID = "createTalkSession"
								r.pathPattern = "/api/talksessions"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							origElem := elem
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "talkSessionId"
							// Match until "/"
							idx := strings.IndexByte(elem, '/')
							if idx < 0 {
								idx = len(elem)
							}
							args[0] = elem[:idx]
							elem = elem[idx:]

							if len(elem) == 0 {
								switch method {
								case "GET":
									r.name = "GetTalkSessionDetail"
									r.summary = "トークセッションの詳細"
									r.operationID = "getTalkSessionDetail"
									r.pathPattern = "/api/talksessions/{talkSessionId}"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}
							switch elem[0] {
							case '/': // Prefix: "/opinions"
								origElem := elem
								if l := len("/opinions"); len(elem) >= l && elem[0:l] == "/opinions" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									switch method {
									case "POST":
										r.name = "PostOpinionPost"
										r.summary = "セッションに対して意見投稿"
										r.operationID = "postOpinionPost"
										r.pathPattern = "/api/talksessions/{talkSessionID}/opinions"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}
								switch elem[0] {
								case '/': // Prefix: "/"
									origElem := elem
									if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
										elem = elem[l:]
									} else {
										break
									}

									// Param: "opinionID"
									// Match until "/"
									idx := strings.IndexByte(elem, '/')
									if idx < 0 {
										idx = len(elem)
									}
									args[1] = elem[:idx]
									elem = elem[idx:]

									if len(elem) == 0 {
										break
									}
									switch elem[0] {
									case '/': // Prefix: "/intentions"
										origElem := elem
										if l := len("/intentions"); len(elem) >= l && elem[0:l] == "/intentions" {
											elem = elem[l:]
										} else {
											break
										}

										if len(elem) == 0 {
											// Leaf node.
											switch method {
											case "POST":
												r.name = "Intention"
												r.summary = "意思表明API"
												r.operationID = "Intention"
												r.pathPattern = "/api/talksessions/{talkSessionID}/opinions/{opinionID}/intentions"
												r.args = args
												r.count = 2
												return r, true
											default:
												return
											}
										}

										elem = origElem
									}

									elem = origElem
								}

								elem = origElem
							}

							elem = origElem
						}

						elem = origElem
					}

					elem = origElem
				case 'u': // Prefix: "user"
					origElem := elem
					if l := len("user"); len(elem) >= l && elem[0:l] == "user" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							r.name = "GetUserProfile"
							r.summary = "ユーザー情報の取得"
							r.operationID = "getUserProfile"
							r.pathPattern = "/api/user"
							r.args = args
							r.count = 0
							return r, true
						case "PUT":
							r.name = "EditUserProfile"
							r.summary = "ユーザー情報の変更"
							r.operationID = "editUserProfile"
							r.pathPattern = "/api/user"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
					switch elem[0] {
					case '/': // Prefix: "/register"
						origElem := elem
						if l := len("/register"); len(elem) >= l && elem[0:l] == "/register" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "POST":
								r.name = "RegisterUser"
								r.summary = "ユーザー作成"
								r.operationID = "registerUser"
								r.pathPattern = "/api/user/register"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'u': // Prefix: "uth/"
				origElem := elem
				if l := len("uth/"); len(elem) >= l && elem[0:l] == "uth/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "provider"
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
					origElem := elem
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'c': // Prefix: "callback"
						origElem := elem
						if l := len("callback"); len(elem) >= l && elem[0:l] == "callback" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "GET":
								r.name = "OAuthCallback"
								r.summary = "Auth Callback"
								r.operationID = "oauth_callback"
								r.pathPattern = "/auth/{provider}/callback"
								r.args = args
								r.count = 1
								return r, true
							default:
								return
							}
						}

						elem = origElem
					case 'l': // Prefix: "login"
						origElem := elem
						if l := len("login"); len(elem) >= l && elem[0:l] == "login" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "GET":
								r.name = "AuthLogin"
								r.summary = "OAuthログイン"
								r.operationID = "auth_login"
								r.pathPattern = "/auth/{provider}/login"
								r.args = args
								r.count = 1
								return r, true
							default:
								return
							}
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			}

			elem = origElem
		}
	}
	return r, false
}
