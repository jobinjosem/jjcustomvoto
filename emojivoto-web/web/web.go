package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/jobinjosem/jjcustomvoto/pkg/api"

	pb "github.com/jobinjosem/jjcustomvoto/emojivoto-web/gen/proto"
	// "github.com/jobinjosem/jjcustomvoto/pkg/api"
	// _ "github.com/jobinjosem/jjcustomvoto/pkg/api/docs"
	// "github.com/swaggo/swag"
	"go.opencensus.io/plugin/ochttp"
	// "github.com/jobinjosem/jjcustomvoto/pkg/grpc"
	// "github.com/jobinjosem/jjcustomvoto/pkg/signals"
	// "github.com/jobinjosem/jjcustomvoto/pkg/version"
	// httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	emojiServiceClient  pb.EmojiServiceClient
	votingServiceClient pb.VotingServiceClient
	indexBundle         string
	webpackDevServer    string
	messageOfTheDay     string
}

func (s *Server) listEmojiHandler(w http.ResponseWriter, r *http.Request) {
	serviceResponse, err := s.emojiServiceClient.ListAll(r.Context(), &pb.ListAllEmojiRequest{})
	if err != nil {
		WriteError(err, w, r, http.StatusInternalServerError, true)
		return
	}

	list := make([]map[string]string, 0)
	for _, e := range serviceResponse.List {
		list = append(list, map[string]string{
			"shortcode": e.Shortcode,
			"unicode":   e.Unicode,
		})
	}

	err = writeJsonBody(w, http.StatusOK, list)

	if err != nil {
		WriteError(err, w, r, http.StatusInternalServerError, true)
	}
}

func (s *Server) leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	results, err := s.votingServiceClient.Results(r.Context(), &pb.ResultsRequest{})

	if err != nil {
		WriteError(err, w, r, http.StatusInternalServerError, true)
		return
	}

	representations := make([]map[string]string, 0)
	for _, result := range results.Results {
		findByShortcodeRequest := &pb.FindByShortcodeRequest{
			Shortcode: result.Shortcode,
		}

		findByShortcodeResponse, err := s.emojiServiceClient.FindByShortcode(r.Context(), findByShortcodeRequest)

		if err != nil {
			WriteError(err, w, r, http.StatusInternalServerError, true)
			return
		}

		emoji := findByShortcodeResponse.Emoji
		representation := make(map[string]string)
		representation["votes"] = strconv.Itoa(int(result.Votes))
		representation["unicode"] = emoji.Unicode
		representation["shortcode"] = emoji.Shortcode

		representations = append(representations, representation)
	}

	err = writeJsonBody(w, http.StatusOK, representations)

	if err != nil {
		WriteError(err, w, r, http.StatusInternalServerError, true)
	}
}

func (s *Server) voteEmojiHandler(w http.ResponseWriter, r *http.Request) {
	emojiShortcode := r.FormValue("choice")
	if emojiShortcode == "" {
		error := errors.New(fmt.Sprintf("Emoji choice [%s] is mandatory", emojiShortcode))
		WriteError(error, w, r, http.StatusBadRequest, true)
		return
	}

	request := &pb.FindByShortcodeRequest{
		Shortcode: emojiShortcode,
	}
	response, err := s.emojiServiceClient.FindByShortcode(r.Context(), request)
	if err != nil {
		WriteError(err, w, r, http.StatusInternalServerError, true)
		return
	}

	if response.Emoji == nil {
		err = errors.New(fmt.Sprintf("Choosen emoji shortcode [%s] doesnt exist", emojiShortcode))
		WriteError(err, w, r, http.StatusBadRequest, true)
		return
	}

	voteRequest := &pb.VoteRequest{}
	switch emojiShortcode {
	case ":poop:":
		_, err = s.votingServiceClient.VotePoop(r.Context(), voteRequest)
	case ":joy:":
		_, err = s.votingServiceClient.VoteJoy(r.Context(), voteRequest)
	case ":sunglasses:":
		_, err = s.votingServiceClient.VoteSunglasses(r.Context(), voteRequest)
	case ":relaxed:":
		_, err = s.votingServiceClient.VoteRelaxed(r.Context(), voteRequest)
	case ":stuck_out_tongue_winking_eye:":
		_, err = s.votingServiceClient.VoteStuckOutTongueWinkingEye(r.Context(), voteRequest)
	case ":money_mouth_face:":
		_, err = s.votingServiceClient.VoteMoneyMouthFace(r.Context(), voteRequest)
	case ":flushed:":
		_, err = s.votingServiceClient.VoteFlushed(r.Context(), voteRequest)
	case ":mask:":
		_, err = s.votingServiceClient.VoteMask(r.Context(), voteRequest)
	case ":nerd_face:":
		_, err = s.votingServiceClient.VoteNerdFace(r.Context(), voteRequest)
	case ":ghost:":
		_, err = s.votingServiceClient.VoteGhost(r.Context(), voteRequest)
	case ":skull_and_crossbones:":
		_, err = s.votingServiceClient.VoteSkullAndCrossbones(r.Context(), voteRequest)
	case ":heart_eyes_cat:":
		_, err = s.votingServiceClient.VoteHeartEyesCat(r.Context(), voteRequest)
	case ":hear_no_evil:":
		_, err = s.votingServiceClient.VoteHearNoEvil(r.Context(), voteRequest)
	case ":see_no_evil:":
		_, err = s.votingServiceClient.VoteSeeNoEvil(r.Context(), voteRequest)
	case ":speak_no_evil:":
		_, err = s.votingServiceClient.VoteSpeakNoEvil(r.Context(), voteRequest)
	case ":boy:":
		_, err = s.votingServiceClient.VoteBoy(r.Context(), voteRequest)
	case ":girl:":
		_, err = s.votingServiceClient.VoteGirl(r.Context(), voteRequest)
	case ":man:":
		_, err = s.votingServiceClient.VoteMan(r.Context(), voteRequest)
	case ":woman:":
		_, err = s.votingServiceClient.VoteWoman(r.Context(), voteRequest)
	case ":older_man:":
		_, err = s.votingServiceClient.VoteOlderMan(r.Context(), voteRequest)
	case ":policeman:":
		_, err = s.votingServiceClient.VotePoliceman(r.Context(), voteRequest)
	case ":guardsman:":
		_, err = s.votingServiceClient.VoteGuardsman(r.Context(), voteRequest)
	case ":construction_worker_man:":
		_, err = s.votingServiceClient.VoteConstructionWorkerMan(r.Context(), voteRequest)
	case ":prince:":
		_, err = s.votingServiceClient.VotePrince(r.Context(), voteRequest)
	case ":princess:":
		_, err = s.votingServiceClient.VotePrincess(r.Context(), voteRequest)
	case ":man_in_tuxedo:":
		_, err = s.votingServiceClient.VoteManInTuxedo(r.Context(), voteRequest)
	case ":bride_with_veil:":
		_, err = s.votingServiceClient.VoteBrideWithVeil(r.Context(), voteRequest)
	case ":mrs_claus:":
		_, err = s.votingServiceClient.VoteMrsClaus(r.Context(), voteRequest)
	case ":santa:":
		_, err = s.votingServiceClient.VoteSanta(r.Context(), voteRequest)
	case ":turkey:":
		_, err = s.votingServiceClient.VoteTurkey(r.Context(), voteRequest)
	case ":rabbit:":
		_, err = s.votingServiceClient.VoteRabbit(r.Context(), voteRequest)
	case ":no_good_woman:":
		_, err = s.votingServiceClient.VoteNoGoodWoman(r.Context(), voteRequest)
	case ":ok_woman:":
		_, err = s.votingServiceClient.VoteOkWoman(r.Context(), voteRequest)
	case ":raising_hand_woman:":
		_, err = s.votingServiceClient.VoteRaisingHandWoman(r.Context(), voteRequest)
	case ":bowing_man:":
		_, err = s.votingServiceClient.VoteBowingMan(r.Context(), voteRequest)
	case ":man_facepalming:":
		_, err = s.votingServiceClient.VoteManFacepalming(r.Context(), voteRequest)
	case ":woman_shrugging:":
		_, err = s.votingServiceClient.VoteWomanShrugging(r.Context(), voteRequest)
	case ":massage_woman:":
		_, err = s.votingServiceClient.VoteMassageWoman(r.Context(), voteRequest)
	case ":walking_man:":
		_, err = s.votingServiceClient.VoteWalkingMan(r.Context(), voteRequest)
	case ":running_man:":
		_, err = s.votingServiceClient.VoteRunningMan(r.Context(), voteRequest)
	case ":dancer:":
		_, err = s.votingServiceClient.VoteDancer(r.Context(), voteRequest)
	case ":man_dancing:":
		_, err = s.votingServiceClient.VoteManDancing(r.Context(), voteRequest)
	case ":dancing_women:":
		_, err = s.votingServiceClient.VoteDancingWomen(r.Context(), voteRequest)
	case ":rainbow:":
		_, err = s.votingServiceClient.VoteRainbow(r.Context(), voteRequest)
	case ":skier:":
		_, err = s.votingServiceClient.VoteSkier(r.Context(), voteRequest)
	case ":golfing_man:":
		_, err = s.votingServiceClient.VoteGolfingMan(r.Context(), voteRequest)
	case ":surfing_man:":
		_, err = s.votingServiceClient.VoteSurfingMan(r.Context(), voteRequest)
	case ":basketball_man:":
		_, err = s.votingServiceClient.VoteBasketballMan(r.Context(), voteRequest)
	case ":biking_man:":
		_, err = s.votingServiceClient.VoteBikingMan(r.Context(), voteRequest)
	case ":point_up_2:":
		_, err = s.votingServiceClient.VotePointUp2(r.Context(), voteRequest)
	case ":vulcan_salute:":
		_, err = s.votingServiceClient.VoteVulcanSalute(r.Context(), voteRequest)
	case ":metal:":
		_, err = s.votingServiceClient.VoteMetal(r.Context(), voteRequest)
	case ":call_me_hand:":
		_, err = s.votingServiceClient.VoteCallMeHand(r.Context(), voteRequest)
	case ":thumbsup:":
		_, err = s.votingServiceClient.VoteThumbsup(r.Context(), voteRequest)
	case ":wave:":
		_, err = s.votingServiceClient.VoteWave(r.Context(), voteRequest)
	case ":clap:":
		_, err = s.votingServiceClient.VoteClap(r.Context(), voteRequest)
	case ":raised_hands:":
		_, err = s.votingServiceClient.VoteRaisedHands(r.Context(), voteRequest)
	case ":pray:":
		_, err = s.votingServiceClient.VotePray(r.Context(), voteRequest)
	case ":dog:":
		_, err = s.votingServiceClient.VoteDog(r.Context(), voteRequest)
	case ":cat2:":
		_, err = s.votingServiceClient.VoteCat2(r.Context(), voteRequest)
	case ":pig:":
		_, err = s.votingServiceClient.VotePig(r.Context(), voteRequest)
	case ":hatching_chick:":
		_, err = s.votingServiceClient.VoteHatchingChick(r.Context(), voteRequest)
	case ":snail:":
		_, err = s.votingServiceClient.VoteSnail(r.Context(), voteRequest)
	case ":bacon:":
		_, err = s.votingServiceClient.VoteBacon(r.Context(), voteRequest)
	case ":pizza:":
		_, err = s.votingServiceClient.VotePizza(r.Context(), voteRequest)
	case ":taco:":
		_, err = s.votingServiceClient.VoteTaco(r.Context(), voteRequest)
	case ":burrito:":
		_, err = s.votingServiceClient.VoteBurrito(r.Context(), voteRequest)
	case ":ramen:":
		_, err = s.votingServiceClient.VoteRamen(r.Context(), voteRequest)
	case ":doughnut:":
		_, err = s.votingServiceClient.VoteDoughnut(r.Context(), voteRequest)
	case ":champagne:":
		_, err = s.votingServiceClient.VoteChampagne(r.Context(), voteRequest)
	case ":tropical_drink:":
		_, err = s.votingServiceClient.VoteTropicalDrink(r.Context(), voteRequest)
	case ":beer:":
		_, err = s.votingServiceClient.VoteBeer(r.Context(), voteRequest)
	case ":tumbler_glass:":
		_, err = s.votingServiceClient.VoteTumblerGlass(r.Context(), voteRequest)
	case ":world_map:":
		_, err = s.votingServiceClient.VoteWorldMap(r.Context(), voteRequest)
	case ":beach_umbrella:":
		_, err = s.votingServiceClient.VoteBeachUmbrella(r.Context(), voteRequest)
	case ":mountain_snow:":
		_, err = s.votingServiceClient.VoteMountainSnow(r.Context(), voteRequest)
	case ":camping:":
		_, err = s.votingServiceClient.VoteCamping(r.Context(), voteRequest)
	case ":steam_locomotive:":
		_, err = s.votingServiceClient.VoteSteamLocomotive(r.Context(), voteRequest)
	case ":flight_departure:":
		_, err = s.votingServiceClient.VoteFlightDeparture(r.Context(), voteRequest)
	case ":rocket:":
		_, err = s.votingServiceClient.VoteRocket(r.Context(), voteRequest)
	case ":star2:":
		_, err = s.votingServiceClient.VoteStar2(r.Context(), voteRequest)
	case ":sun_behind_small_cloud:":
		_, err = s.votingServiceClient.VoteSunBehindSmallCloud(r.Context(), voteRequest)
	case ":cloud_with_rain:":
		_, err = s.votingServiceClient.VoteCloudWithRain(r.Context(), voteRequest)
	case ":fire:":
		_, err = s.votingServiceClient.VoteFire(r.Context(), voteRequest)
	case ":jack_o_lantern:":
		_, err = s.votingServiceClient.VoteJackOLantern(r.Context(), voteRequest)
	case ":balloon:":
		_, err = s.votingServiceClient.VoteBalloon(r.Context(), voteRequest)
	case ":tada:":
		_, err = s.votingServiceClient.VoteTada(r.Context(), voteRequest)
	case ":trophy:":
		_, err = s.votingServiceClient.VoteTrophy(r.Context(), voteRequest)
	case ":iphone:":
		_, err = s.votingServiceClient.VoteIphone(r.Context(), voteRequest)
	case ":pager:":
		_, err = s.votingServiceClient.VotePager(r.Context(), voteRequest)
	case ":fax:":
		_, err = s.votingServiceClient.VoteFax(r.Context(), voteRequest)
	case ":bulb:":
		_, err = s.votingServiceClient.VoteBulb(r.Context(), voteRequest)
	case ":money_with_wings:":
		_, err = s.votingServiceClient.VoteMoneyWithWings(r.Context(), voteRequest)
	case ":crystal_ball:":
		_, err = s.votingServiceClient.VoteCrystalBall(r.Context(), voteRequest)
	case ":underage:":
		_, err = s.votingServiceClient.VoteUnderage(r.Context(), voteRequest)
	case ":interrobang:":
		_, err = s.votingServiceClient.VoteInterrobang(r.Context(), voteRequest)
	case ":100:":
		_, err = s.votingServiceClient.Vote100(r.Context(), voteRequest)
	case ":checkered_flag:":
		_, err = s.votingServiceClient.VoteCheckeredFlag(r.Context(), voteRequest)
	case ":crossed_swords:":
		_, err = s.votingServiceClient.VoteCrossedSwords(r.Context(), voteRequest)
	case ":floppy_disk:":
		_, err = s.votingServiceClient.VoteFloppyDisk(r.Context(), voteRequest)
	}
	if err != nil {
		WriteError(err, w, r, http.StatusInternalServerError, true)
		return
	}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	indexTemplate := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Emoji Vote</title>
			<link rel="icon" href="/img/favicon.ico">
			<!-- Global site tag (gtag.js) - Google Analytics -->
			<script async src="https://www.googletagmanager.com/gtag/js?id=UA-60040560-4"></script>
			<script>
			  window.dataLayer = window.dataLayer || [];
			  function gtag(){dataLayer.push(arguments);}
			  gtag('js', new Date());
			  gtag('config', 'UA-60040560-4');
			</script>
		</head>
		<body>
			<div id="motd" class="motd">%s</div>
			<div id="main" class="main"></div>
		</body>
		{{ if ne . ""}}
			<script type="text/javascript" src="{{ . }}/dist/index_bundle.js" async></script>
		{{else}}
			<script type="text/javascript" src="/js" async></script>
		{{end}}
	</html>`, s.messageOfTheDay)
	t, err := template.New("indexTemplate").Parse(indexTemplate)
	if err != nil {
		panic(err)
	}
	t.Execute(w, s.webpackDevServer)
}

func (s *Server) jsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	f, err := ioutil.ReadFile(s.indexBundle)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(f))
}

func (s *Server) faviconHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			WriteError(fmt.Errorf("%v", err), w, r, http.StatusInternalServerError, true)
		}
	}()

	http.ServeFile(w, r, "./web/favicon.ico")
}

func writeJsonBody(w http.ResponseWriter, status int, body interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

func WriteError(err error, w http.ResponseWriter, r *http.Request, status int, debug bool) {
	logMessage := fmt.Sprintf("Error serving request [%v]: %v", r, err)

	if debug {
		logMessage += fmt.Sprintf("\nRequest Headers: %+v", r.Header)
		logMessage += fmt.Sprintf("\nRequest Body: %+v", r.Body)
	}

	log.Printf(logMessage)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	errorMessage := make(map[string]interface{})
	errorMessage["error"] = fmt.Sprintf("%v", err)

	if debug {
		errorMessage["method"] = r.Method
		errorMessage["url"] = r.URL.String()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			body = []byte{}
		}
		errorMessage["request_body"] = string(body)
	}

	json.NewEncoder(w).Encode(errorMessage)
}

func handle(path string, h func(w http.ResponseWriter, r *http.Request)) {
	http.Handle(path, &ochttp.Handler{
		Handler: http.HandlerFunc(h),
	})
}

func StartServer(webPort, webpackDevServer, indexBundle string, emojiServiceClient pb.EmojiServiceClient, votingClient pb.VotingServiceClient) {

	motd := os.Getenv("MESSAGE_OF_THE_DAY")
	Server := &Server{
		emojiServiceClient:  emojiServiceClient,
		votingServiceClient: votingClient,
		indexBundle:         indexBundle,
		webpackDevServer:    webpackDevServer,
		messageOfTheDay:     motd,
	}
	// API := &Api{
    //    router
	// }

	log.Printf("Starting web server on WEB_PORT=[%s] and MESSAGE_OF_THE_DAY=[%s]", webPort, motd)
	handle("/", Server.indexHandler)
	handle("/leaderboard", Server.indexHandler)
	handle("/js", Server.jsHandler)
	handle("/img/favicon.ico", Server.faviconHandler)
	handle("/api/list", Server.listEmojiHandler)
	handle("/api/vote", Server.voteEmojiHandler)
	handle("/api/leaderboard", Server.leaderboardHandler)
	// handle("/env", api.NewMockServer().EnvHandler)
	handle("/version", Api.VersionHandler)
	// handle("/info", api.NewMockServer().InfoHandler)
	// Server.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
    //     httpSwagger.URL("/swagger/doc.json"),
    // ))
    // Server.router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
    //     doc, err := swag.ReadDoc()
    //     if err != nil {
    //         Server.logger.Error("swagger error", zap.Error(err), zap.String("path", "/swagger.json"))
    //     }
    //     w.Write([]byte(doc))
    // })
	// TODO: make static assets dir configurable
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	err := http.ListenAndServe(fmt.Sprintf(":%s", webPort), nil)
	if err != nil {
		panic(err)
	}
}
