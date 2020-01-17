(function (window, undifined) {
    var wsUri = "ws://127.0.0.1:7777/ws",
        player = [],
        team = [],
        teamNum,
        time,
        sdc,
        info,
        pay_box,
        vip,
        input_1,
        limit,
        win,
        text_2,
        choice

    team[1] = document.querySelector("#team_1")
    team[2] = document.querySelector("#team_2")
    team[3] = document.querySelector("#team_3")
    player[1] = document.querySelector("#player_1")
    player[2] = document.querySelector("#player_2")
    player[3] = document.querySelector("#player_3")
    sdc = document.querySelector("#sdc")
    time = document.querySelector("#time")
    info = document.querySelector("#info")
    pay_box = document.querySelector("#pay_box")
    input_1 = document.querySelector("#input_1")
    limit = document.querySelector("#limit")
    text_2 = document.querySelector("#text_2")
    choice = document.querySelectorAll(".choice")

    function init() {
        testWebSocket();
        timer()
        setTimeout(_ => {
            document.querySelector('#load').style.opacity = 0
            setTimeout(_ => {
                document.querySelector('#load').style.display = 'none'
            }, 500)
        }, 1000)
    }

    function testWebSocket() {
        websocket = new WebSocket(wsUri);
        websocket.onopen = function () {
            login()
        };
        websocket.onclose = function () { };
        websocket.onmessage = function (evt) {
            console.log(evt.data)
            var data = JSON.parse(evt.data)
            if (data.hasOwnProperty("TortInfo")) {
                TortInfo(data["TortInfo"])
            }

            if (data.hasOwnProperty("LoginResult")) {
                LoginResult(data["LoginResult"])
            }

            if (data.hasOwnProperty("AwardResult")) {
                AwardResult(data["AwardResult"])
            }
        };
        websocket.onerror = function () { };
    }

    window.closePrompt = function (item, key) {
        setTimeout(_ => {
            document.querySelector(`#prompt_${item}`).style.display = 'none'
        }, 200)
        document.querySelector(`#background_${item}`).style.opacity = 0
        document.querySelector(`#background_${item}_cover`).style.opacity = 0
        document.querySelector(`#box_${item}`).style.height = 0
        if (teamNum && key) {
            document.querySelector(`#message_${teamNum}`).innerHTML = `真是遗憾`
            document.querySelector(`#message_${teamNum}`).style.display = `block`
            setTimeout(_ => {
                document.querySelector(`#message_${teamNum}`).style.display = `none`
            }, 1000)
        }
    }

    //游戏开始
    function GameStart(data) {
        player[1].style.bottom = `${data[0] * 100}%`
        player[2].style.bottom = `${data[1] * 100}%`
        player[3].style.bottom = `${data[2] * 100}%`
    }

    //显示队伍信息及收益
    function TortInfo(data) {
        var text = ``
        data["isPutSDC"] ? text = `选择时间` : ''
        data["isGameEnd"] ? gameEnd(data["gameResult"]) : ''
        timer(data["time"])
        TortTeam(data["teamInfo"])
        GameStart(data["gameResult"])
    }

    function gameEnd(data) {
        let i = 0
        data[0] > data[1] ? i = 0 : i = 1
        data[i] > data[2] ? '' : i = 2
        if (win) return
        win = document.querySelector(`#player_${i + 1}_bt`)
        win.style.display == 'block' ? '' : win.style.display = 'block'
    }

    function TortTeam(data) {
        var ss = data[0][1] + data[1][1] + data[2][1]
        data[0][2] = (0.9 * (ss - data[0][1])) / (data[0][1] || 1) + 1
        data[1][2] = (0.9 * (ss - data[1][1])) / (data[1][1] || 1) + 1
        data[2][2] = (0.9 * (ss - data[2][1])) / (data[2][1] || 1) + 1
        teamInfo(1, data[0])
        teamInfo(2, data[1])
        teamInfo(3, data[2])
    }

    function teamInfo(item, data) {
        team[item].innerHTML = `<div>${data[0]}</div><div>${data[1]}</div><div>${parseInt(data[2] * 100) / 100}</div>`
    }

    function LoginResult(data) {
        // console.log(data)
        if (data.code == 0) {
            sdc.innerHTML = data.sdc ? data.sdc : 0
            vip = data.vip
            limit.innerHTML = `0~${(vip + 1) * 10}`
            if (data.team) {
                choice[data.team - 1].style.display = 'block'
            }
            if (teamNum) {
                document.querySelector(`#message_${teamNum}`).innerHTML = `感觉能赢`
                document.querySelector(`#message_${teamNum}`).style.display = `block`
                setTimeout(_ => {
                    document.querySelector(`#message_${teamNum}`).style.display = `none`
                }, 1000)
                closePrompt(1)
                return
            }
            closePrompt(1, true)
            return
        }
        sdc.innerHTML = 0
        dealErr(data)
    }

    function AwardResult(data) {
        if (data.code == 0) {
            text_2.innerHTML = `恭喜您的队伍获得胜利,您将获得${data.sdc}神果奖励`
            return
        }
        text_2.innerHTML = data.msg
        document.querySelector(`#prompt_2`).style.display = 'block'
        setTimeout(_ => {
            document.querySelector(`#background_2`).style.opacity = 1
            document.querySelector(`#background_2_cover`).style.opacity = 1
            document.querySelector(`#box_2`).style.height = '140px'
        }, 20)
    }

    function doSend(message) {
        websocket.send(message);
    }

    function dealErr(key) {
        switch (key.code) {
            case -1:
                limit.innerHTML = `${key.msg}`
                break;
            case -2:
                limit.innerHTML = `${key.msg}`
                break
            default:
                limit.innerHTML = `${key.msg}`
                break;
        }
    }

    function timer(t) {
        h = new Date(t).getHours()
        m = new Date(t).getMinutes()
        s = new Date(t).getSeconds()
        // console.log(h,m,s)
        if (h == 19 && m >= 58) {
            pay_box.style.display == 'none' ? '' : pay_box.style.display = 'none'
            time.innerHTML = `0${59 - m}:${s > 49 ? '0' + (59 - s) : 59 - s}`
            info.innerHTML = `准备时间`
        } else if (h == 20 && m <= 1) {
            time.innerHTML = `00:${s > 49 ? '0' + (59 - s) : 59 - s}`
            info.innerHTML = `游戏开始`
        } else if (h >= 20 && m > 1 && m < 5) {
            time.innerHTML = `0${4 - m}:${s > 49 ? '0' + (59 - s) : 59 - s}`
            info.innerHTML = `游戏结算`
        } else if (h >= 20 && m >= 5) {
            time.innerHTML = `0${23 - h}:${m > 49 ? '0' + (59 - m) : 59 - m}`
            info.innerHTML = `领取奖励`
        } else {
            if (win) {
                win.style.display == 'none' ? '' : win.style.display = 'none'
                win = null
            }
            info.innerHTML = `当前时间`
            pay_box.style.display == 'flex' ? '' : pay_box.style.display = 'flex'
            time.innerHTML = `${h < 10 ? '0' + h : h}:${m < 10 ? '0' + m : m}`
        }
    }

    window.addEventListener("load", init, false);

    function login() {
        var query = window.location.search.substring(1);
        var arr = query.split("&")
        var msg = {}
        for (let i = 0; i < arr.length; i++) {
            var pair = arr[i].split("=");
            msg[pair[0]] = pair[1]
        }
        doSend("login-" + JSON.stringify(msg));
    }

    window.openJion = function (item) {
        document.querySelector(`#prompt_1`).style.display = 'block'
        setTimeout(_ => {
            document.querySelector(`#background_1`).style.opacity = 1
            document.querySelector(`#background_1_cover`).style.opacity = 1
            document.querySelector(`#box_1`).style.height = '140px'
        }, 20)
        teamNum = item
    }

    window.jionTeam = function () {
        var query = window.location.search.substring(1);
        var arr = query.split("&")
        var msg = {}
        for (let i = 0; i < arr.length; i++) {
            var pair = arr[i].split("=");
            msg[pair[0]] = pair[1]
        }
        let sdc = input_1.value
        if (!sdc) {
            return limit.innerHTML = `输入错误`
        }
        if (!team) {
            return limit.innerHTML = `网络崩溃,请稍后重试`
        }
        msg["sdc"] = parseFloat(sdc)
        msg["team"] = teamNum
        doSend("jion--" + JSON.stringify(msg));
    }

    window.award = function () {
        var query = window.location.search.substring(1);
        var arr = query.split("&")
        var msg = {}
        for (let i = 0; i < arr.length; i++) {
            var pair = arr[i].split("=");
            msg[pair[0]] = pair[1]
        }
        doSend("award-" + JSON.stringify(msg))
    }

    window.closeBtnClick = function () {
        websocket.close();
    }
})(window)