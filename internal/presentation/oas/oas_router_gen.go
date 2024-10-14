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
			case 'a': // Prefix: "auth/"
				origElem := elem
				if l := len("auth/"); len(elem) >= l && elem[0:l] == "auth/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'r': // Prefix: "revoke"
					origElem := elem
					if l := len("revoke"); len(elem) >= l && elem[0:l] == "revoke" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleOAuthRevokeRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}

					elem = origElem
				case 't': // Prefix: "token/info"
					origElem := elem
					if l := len("token/info"); len(elem) >= l && elem[0:l] == "token/info" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleOAuthTokenInfoRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}

					elem = origElem
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
								s.handleAuthorizeRequest([1]string{
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
			case 'o': // Prefix: "opinions/histories"
				origElem := elem
				if l := len("opinions/histories"); len(elem) >= l && elem[0:l] == "opinions/histories" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "GET":
						s.handleOpinionsHistoryRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET")
					}

					return
				}

				elem = origElem
			case 's': // Prefix: "sessions/histories"
				origElem := elem
				if l := len("sessions/histories"); len(elem) >= l && elem[0:l] == "sessions/histories" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "GET":
						s.handleSessionsHistoryRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET")
					}

					return
				}

				elem = origElem
			case 't': // Prefix: "t"
				origElem := elem
				if l := len("t"); len(elem) >= l && elem[0:l] == "t" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'a': // Prefix: "alksession"
					origElem := elem
					if l := len("alksession"); len(elem) >= l && elem[0:l] == "alksession" {
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
							case 'o': // Prefix: "opinions/"
								origElem := elem
								if l := len("opinions/"); len(elem) >= l && elem[0:l] == "opinions/" {
									elem = elem[l:]
								} else {
									break
								}

								// Param: "opinionID"
								// Leaf parameter
								args[1] = elem
								elem = ""

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetOpinionDetailRequest([2]string{
											args[0],
											args[1],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}

								elem = origElem
							case 's': // Prefix: "swipe_opinions"
								origElem := elem
								if l := len("swipe_opinions"); len(elem) >= l && elem[0:l] == "swipe_opinions" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleSwipeOpinionsRequest([1]string{
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
								s.handleGetTalkSessionListRequest([0]string{}, elemIsEscaped, w, r)
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
							case '/': // Prefix: "/opinion"
								origElem := elem
								if l := len("/opinion"); len(elem) >= l && elem[0:l] == "/opinion" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									switch r.Method {
									case "GET":
										s.handleGetTopOpinionsRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}
								switch elem[0] {
								case 's': // Prefix: "s"
									origElem := elem
									if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
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
											case 'r': // Prefix: "replies"
												origElem := elem
												if l := len("replies"); len(elem) >= l && elem[0:l] == "replies" {
													elem = elem[l:]
												} else {
													break
												}

												if len(elem) == 0 {
													// Leaf node.
													switch r.Method {
													case "GET":
														s.handleOpinionCommentsRequest([2]string{
															args[0],
															args[1],
														}, elemIsEscaped, w, r)
													default:
														s.notAllowed(w, r, "GET")
													}

													return
												}

												elem = origElem
											case 'v': // Prefix: "votes"
												origElem := elem
												if l := len("votes"); len(elem) >= l && elem[0:l] == "votes" {
													elem = elem[l:]
												} else {
													break
												}

												if len(elem) == 0 {
													// Leaf node.
													switch r.Method {
													case "POST":
														s.handleVoteRequest([2]string{
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
						}

						elem = origElem
					}

					elem = origElem
				case 'e': // Prefix: "est"
					origElem := elem
					if l := len("est"); len(elem) >= l && elem[0:l] == "est" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleTestRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
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
					// Leaf node.
					switch r.Method {
					case "GET":
						s.handleGetUserInfoRequest([0]string{}, elemIsEscaped, w, r)
					case "POST":
						s.handleRegisterUserRequest([0]string{}, elemIsEscaped, w, r)
					case "PUT":
						s.handleEditUserProfileRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET,POST,PUT")
					}

					return
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
			case 'a': // Prefix: "auth/"
				origElem := elem
				if l := len("auth/"); len(elem) >= l && elem[0:l] == "auth/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'r': // Prefix: "revoke"
					origElem := elem
					if l := len("revoke"); len(elem) >= l && elem[0:l] == "revoke" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "POST":
							r.name = "OAuthRevoke"
							r.summary = "アクセストークンを失効"
							r.operationID = "oauth_revoke"
							r.pathPattern = "/auth/revoke"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				case 't': // Prefix: "token/info"
					origElem := elem
					if l := len("token/info"); len(elem) >= l && elem[0:l] == "token/info" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "GET":
							r.name = "OAuthTokenInfo"
							r.summary = "JWTの内容を返してくれる"
							r.operationID = "oauth_token_info"
							r.pathPattern = "/auth/token/info"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
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
								r.name = "Authorize"
								r.summary = "OAuthログイン"
								r.operationID = "authorize"
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
			case 'o': // Prefix: "opinions/histories"
				origElem := elem
				if l := len("opinions/histories"); len(elem) >= l && elem[0:l] == "opinions/histories" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "GET":
						r.name = "OpinionsHistory"
						r.summary = "今までに投稿した異見"
						r.operationID = "opinionsHistory"
						r.pathPattern = "/opinions/histories"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}

				elem = origElem
			case 's': // Prefix: "sessions/histories"
				origElem := elem
				if l := len("sessions/histories"); len(elem) >= l && elem[0:l] == "sessions/histories" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "GET":
						r.name = "SessionsHistory"
						r.summary = "リアクション済みのセッション一覧"
						r.operationID = "sessionsHistory"
						r.pathPattern = "/sessions/histories"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}

				elem = origElem
			case 't': // Prefix: "t"
				origElem := elem
				if l := len("t"); len(elem) >= l && elem[0:l] == "t" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'a': // Prefix: "alksession"
					origElem := elem
					if l := len("alksession"); len(elem) >= l && elem[0:l] == "alksession" {
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
							case 'o': // Prefix: "opinions/"
								origElem := elem
								if l := len("opinions/"); len(elem) >= l && elem[0:l] == "opinions/" {
									elem = elem[l:]
								} else {
									break
								}

								// Param: "opinionID"
								// Leaf parameter
								args[1] = elem
								elem = ""

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "GET":
										r.name = "GetOpinionDetail"
										r.summary = "意見の詳細"
										r.operationID = "getOpinionDetail"
										r.pathPattern = "/talksession/{talkSessionID}/opinions/{opinionID}"
										r.args = args
										r.count = 2
										return r, true
									default:
										return
									}
								}

								elem = origElem
							case 's': // Prefix: "swipe_opinions"
								origElem := elem
								if l := len("swipe_opinions"); len(elem) >= l && elem[0:l] == "swipe_opinions" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch method {
									case "GET":
										r.name = "SwipeOpinions"
										r.summary = "スワイプ用のエンドポイント"
										r.operationID = "swipe_opinions"
										r.pathPattern = "/talksession/{talkSessionID}/swipe_opinions"
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
								r.name = "GetTalkSessionList"
								r.summary = "トークセッションコレクション"
								r.operationID = "getTalkSessionList"
								r.pathPattern = "/talksessions"
								r.args = args
								r.count = 0
								return r, true
							case "POST":
								r.name = "CreateTalkSession"
								r.summary = "トークセッション作成"
								r.operationID = "createTalkSession"
								r.pathPattern = "/talksessions"
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
									r.summary = "🚧 トークセッションの詳細"
									r.operationID = "getTalkSessionDetail"
									r.pathPattern = "/talksessions/{talkSessionId}"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}
							switch elem[0] {
							case '/': // Prefix: "/opinion"
								origElem := elem
								if l := len("/opinion"); len(elem) >= l && elem[0:l] == "/opinion" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									switch method {
									case "GET":
										r.name = "GetTopOpinions"
										r.summary = "🚧 分析に関する意見"
										r.operationID = "getTopOpinions"
										r.pathPattern = "/talksessions/{talkSessionId}/opinion"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}
								switch elem[0] {
								case 's': // Prefix: "s"
									origElem := elem
									if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
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
											r.pathPattern = "/talksessions/{talkSessionID}/opinions"
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
											case 'r': // Prefix: "replies"
												origElem := elem
												if l := len("replies"); len(elem) >= l && elem[0:l] == "replies" {
													elem = elem[l:]
												} else {
													break
												}

												if len(elem) == 0 {
													// Leaf node.
													switch method {
													case "GET":
														r.name = "OpinionComments"
														r.summary = "意見に対するコメント一覧を返す"
														r.operationID = "opinionComments"
														r.pathPattern = "/talksessions/{talkSessionID}/opinions/{opinionID}/replies"
														r.args = args
														r.count = 2
														return r, true
													default:
														return
													}
												}

												elem = origElem
											case 'v': // Prefix: "votes"
												origElem := elem
												if l := len("votes"); len(elem) >= l && elem[0:l] == "votes" {
													elem = elem[l:]
												} else {
													break
												}

												if len(elem) == 0 {
													// Leaf node.
													switch method {
													case "POST":
														r.name = "Vote"
														r.summary = "意思表明API"
														r.operationID = "vote"
														r.pathPattern = "/talksessions/{talkSessionID}/opinions/{opinionID}/votes"
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
						}

						elem = origElem
					}

					elem = origElem
				case 'e': // Prefix: "est"
					origElem := elem
					if l := len("est"); len(elem) >= l && elem[0:l] == "est" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "GET":
							r.name = "Test"
							r.summary = "OpenAPIテスト用"
							r.operationID = "test"
							r.pathPattern = "/test"
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
			case 'u': // Prefix: "user"
				origElem := elem
				if l := len("user"); len(elem) >= l && elem[0:l] == "user" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "GET":
						r.name = "GetUserInfo"
						r.summary = "ユーザー情報の取得"
						r.operationID = "get_user_info"
						r.pathPattern = "/user"
						r.args = args
						r.count = 0
						return r, true
					case "POST":
						r.name = "RegisterUser"
						r.summary = "ユーザー作成"
						r.operationID = "registerUser"
						r.pathPattern = "/user"
						r.args = args
						r.count = 0
						return r, true
					case "PUT":
						r.name = "EditUserProfile"
						r.summary = "ユーザー情報の変更"
						r.operationID = "editUserProfile"
						r.pathPattern = "/user"
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
	}
	return r, false
}
