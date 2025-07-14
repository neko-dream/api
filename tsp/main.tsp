import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

// Config imports
import "./config/service.tsp";

// Model imports
import "./models/common.tsp";
import "./models/user.tsp";
import "./models/auth.tsp";
import "./models/opinion.tsp";
import "./models/talk-session.tsp";
import "./models/organization.tsp";
import "./models/manage.tsp";

// Route imports
import "./routes/auth.tsp";
import "./routes/user.tsp";
import "./routes/opinion.tsp";
import "./routes/vote.tsp";
import "./routes/talk-session.tsp";
import "./routes/timeline.tsp";
import "./routes/organization.tsp";
import "./routes/policy.tsp";
import "./routes/image.tsp";
import "./routes/test.tsp";
import "./routes/health.tsp";
import "./routes/manage.tsp";
import "./routes/analysis.tsp";

using Http;
using OpenAPI;

@service(#{ title: "Kotohiro API", version: "1.0.0" })
@info(#{ version: "1.0.0" })
@useAuth(CookieAuth)
namespace kotohiro {

}
