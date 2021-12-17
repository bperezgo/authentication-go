const express = require("express");
const path = require("path");
const request = require("request");
const queryString = require("querystring");
const cors = require("cors");
const cookieParser = require("cookie-parser");

const generateRandomString = require("./utils/generateRandomString");
const encodeBasic = require("./utils/encodeBasic");
const scopesArray = require("./utils/scopesArray");

const playlistMocks = require("./utils/mocks/playlist");

const { config } = require("./config");

const app = express();

// static files
app.use("/static", express.static(path.join(__dirname, "public")));

// view engine setup
app.set("views", path.join(__dirname, "views"));
app.set("view engine", "pug");

app.use(cors());
app.use(cookieParser());

// routes
app.get("/", async function (req, res, next) {
  res.render("posts", {
    posts: [
      {
        title: "Guillermo's playlist",
        description:
          "Creatine supplementation is the reference compound for increasing muscular creatine levels; there is variability in this increase, however, with some nonresponders.",
        author: "Guillermo Rodas",
      },
    ],
  });
});

app.get("/login", async (req, res) => {
  const state = generateRandomString(16);
  const query = queryString.stringify({
    // authorization code grant
    response_type: "code",
    client_id: config.spotifyClientId,
    scope: scopesArray.join(" "),
    redirect_uri: config.spotifyRedirectUrl,
    // state to avoid XSS
    state,
  });

  res.cookie("auth_state", state, { httpOnly: true });
  res.redirect(`https://accounts.spotify.com/authorize?${query}`);
});

app.get("/callback", async (req, res, next) => {
  const { code, state } = req.query;
  const { auth_state } = req.cookies;
  if (state === null || state !== auth_state) {
    next(new Error("The state doesn't match"));
  }
  res.clearCookie("auth_state");

  const authOptions = {
    url: "https://accounts.spotify.com/api/token",
    form: {
      code,
      redirect_uri: config.spotifyRedirectUrl,
      grant_type: "authorization_code",
    },
    headers: {
      Authorization: `Basic ${encodeBasic(
        config.spotifyClientId,
        config.spotifyClientSecret
      )}`,
    },
    json: true,
  };

  request.post(authOptions, function (err, res, body) {
    if (err || res.statusCode !== 200) {
      next(new Error("The token is invalid"));
    }

    res.cookie("acces_token", body.access_token, { httpOnly: true });
    res.redirect("/playlist");
  });
});

// server
const server = app.listen(3000, function () {
  console.log(`Listening http://localhost:${server.address().port}`);
});
