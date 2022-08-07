angular.module("logi2").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav", "$http"]
const buttonR = document.getElementById('res');
const buttonErr = document.getElementById('btnerr');
const buttonInf = document.getElementById('btninf');
const buttonDbgs = document.getElementById('btndbgs');
const buttonWar = document.getElementById('btnwar');
const buttonAll = document.getElementById('btnall');
const buttonClr = document.getElementById('changeclr');
const inputform = document.getElementById('search_string');
const textList = document.getElementById('listfile');
const buttonNext = document.getElementById('buttonNext');
const buttonPrev = document.getElementById('buttonPrev');
//
const prevButton = document.getElementById('button_prev');
const nextButton = document.getElementById('button_next');
const clickPageNumber = document.querySelectorAll('.clickPageNumber');

const buttonPage = document.getElementById('page');

const pageNumbers = (total, max, current) => {
    const half = Math.floor(max / 2);
    let to = max;

    if (current + half >= total) {
        to = total;
    } else if (current > half) {
        to = current + half;
    }

    let from = Math.max(to - max, 0);

    return Array.from({ length: Math.min(total, max) }, (_, i) => (i + 1) + from);
}
var countRows = 0
var countWar = 0;
var countErr = 0;
var countInf = 0;
var countDbg = 0;
var countFtl = 0;
var countAll = 0;
var start;
var standartform = "";
var lastItem;
var statusS = "empty";
var typePage = 0;
var numPages;
var ws;
var pageMap;
var Page;

/* 
$(document).ready(function() {
    $('#MyTable').DataTable({
        initComplete: function() {
            this.api().columns().every(function() {
                var column = this;
                var select = $('<select><option value=""></option></select>')
                    .appendTo($(column.footer()).empty())
                    .on('change', function() {
                        var val = $.fn.dataTable.util.escapeRegex(
                            $(this).val()
                        );
                        //to select and search from grid  
                        column
                            .search(val ? '^' + val + '$' : '', true, false)
                            .draw();
                    });

                column.data().unique().sort().each(function(d, j) {
                    select.append('<option value="' + d + '">' + d + '</option>')
                });
            });
        }
    });
});
 */



function isEmpty(str) {
    return (!str || 0 === str.length);
}

buttonR.addEventListener('click', event => {
    setTimeout(
        () => {
            window.location.reload();
        },
        1 * 200
    );
});


inputform.addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        // code for enter
        setTimeout(
            () => {
                Null()
                initWS(lastItem, statusS)
                Null()
            },
            1 * 200
        );
    }
});

/* document.getElementById('pagination11').addEventListener('click', event => {
    setTimeout(
        () => {
            console.log("Work")

        },
        1 * 200
    );
}); */

function editInf() {
    Null()
    initWS(lastItem, "INFO")
    statusS = "INFO"
}

function editErr() {
    Null()
    initWS(lastItem, "ERROR")
    statusS = "ERROR"
}

function editDbgs() {
    Null()
    initWS(lastItem, "DEBUG")
    statusS = "DEBUG"
}

function editWarn() {
    Null()
    initWS(lastItem, "WARNING")
    statusS = "WARNING"
}

function editAll() {
    Null()
    initWS(lastItem, "empty")
    statusS = "empty"
}


function Null() {
    countWar = 0
    countErr = 0
    countInf = 0
    countDbg = 0
    countFtl = 0
    countAll = 0
}

function setBackColor(btn, color) {
    var property = document.getElementById(btn);
    property.style.backgroundColor = color

}

function setFontColor(btn, color) {
    var property = document.getElementById(btn);
    property.style.color = color

}

function quotation(id, text) {
    var q = document.getElementById(id);
    if (q) q.innerHTML = text;
}


function change(identifier, color) {
    identifier.style.background = color;
}

function mainController($rootScope, $scope, $mdSidenav, $http) {

    var vm = this;
    //var lastItem;

    vm.toggleSideNav = function toggleSideNav() {
            $mdSidenav('left').toggle()
        }
        //transmit to server
    vm.init = function init() {
        console.log("In the main controller")
        $scope.showCard = true;
        $http.get('searchproject')
            .then(function(result) {
                $rootScope.search_string = result.data["search_string"]
                console.log("Search :", result.data)
            }, function(result) {
                console.log("Failed to get search")
            })
    }
    vm = this;
    vm.init = function init() {
        console.log("In the main controller")
        $scope.showCard = true;
        $http.get('datestartend')
            .then(function(result) {
                $rootScope.daterange = result.data["daterange"]
                console.log("Search :", result.data)
            }, function(result) {
                console.log("Failed to get search")
            })
    }

    // vm.fontSize = ["10px", "11px", "12px", "14px", "16px", "18px", "20px", "22px", "24px"]
    // $scope.currSize = vm.fontSize[2];


    $scope.open_connection = function(file) {
        var filename = file.replace(/^.*[\\\/]/, '')
        lastItem = null
        lastItem = file;


        /*   console.log(filename) */
        $scope.showCard = false;
        angular.element(document.querySelector("#filename")).html("File: " + filename)

        var container = angular.element(document.querySelector("#container"))

        // var ws;
        if (window.WebSocket === undefined) {
            container.append("Your browser does not support WebSockets");
            return;
        } else {
            // highlight("tbl92")
            ws = initWS(file, "empty");
            Null();
            //ws.close();

        }

        vm.toggleSideNav()
    }

    vm.init();
}


function initWS(file, type) {

    var ws_proto = "ws:"
    if (window.location.protocol === "https:") {
        ws_proto = "wss:"
    }
    console.log(file);
    console.log((ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file)));
    var socket = new WebSocket(ws_proto + "//" + window.location.hostname + ":" + window.location.port + "/ws/" + btoa(file));
    var container = angular.element(document.querySelector("#container"));
    var sysMsgAll = angular.element(document.querySelector("#sysMsgAll"));
    //var count = 0


    container.html("")
    socket.onopen = function(e) {
        strf = file
        if (strf.indexOf("undefined") != 0) {

            // container.append("Start table");

        }
        /*  console.log(file); */
        socket.send(file.replace(/^.*[\\\/]/, ''));

    }

    socket.onmessage = function(e) {
        var loglist
        var map
        str = e.data.trim();
        parser = new DOMParser();
        xmlDoc = parser.parseFromString(str, "text/xml");
        loglist = xmlDoc.getElementsByTagName("loglist");
        map = xmlDoc.getElementsByTagName("Map");
        start = xmlDoc.getElementsByTagName("start");
        countpage = xmlDoc.getElementsByTagName("countpage");

        //console.log("It is str", str)
        k2 = isEmpty(loglist);
        k3 = isEmpty(map);
        k4 = isEmpty(countpage);
        k5 = isEmpty(start)


        if (k5 == false) {
            //console.log("Почалось");
            clearBox("container")
            countWar = 0;
            countErr = 0;
            countInf = 0;
            countDbg = 0;
            countAll = 0;
            var observer = new MutationObserver(function(_mutations, me) {
                // `mutations` is an array of mutations that occurred
                // `me` is the MutationObserver instance
                start = document.getElementById('Foxtrot');
                if (start) {
                    handleCanvas(start);
                    me.disconnect(); // stop observing
                    return;
                }
            });

            // start observing
            observer.observe(document, {
                childList: true,
                subtree: true
            });


        }

        if (k4 == false) {

            $(".pagtest").empty();
            numPages = parseInt(ParseCount(str));

            const paginationButtons = new PaginationButton(numPages, 5);

            paginationButtons.render();

            paginationButtons.onChange(e => {
                console.log("socket send");
                socket.send(e.target.value);
                console.log("Websocket status connection", socket.readyState);
                // console.log("Onmessage connection", socket.readyState);
                console.log(e);
                console.log('-- changed', e.target.value);
            });
        }
        if (k3 == false) {
            //console.log("Websocket status connection onmessage 4", socket.readyState);
            PageMap = str
            str = ParseXmlMap(str)
            quotation("mapContainer", str)

        }
        if (k2 == false) {


            console.log(str)
            str = ParseXml(str, type)
            container.append(str);
            //console.log("Websocket status connection onmessage 5", socket.readyState);
            // socket.close();


            quotation("cntinfo", "INFO:" + countInf);
            quotation("cnterror", "Error:" + countErr);
            quotation("cntwrng", "Warning:" + countWar);
            quotation("ctndbg", "Debug:" + countDbg);
            quotation("cntall", "All:" + countAll);


        } else {
            //:TODO
            //Loading 
            /* if (str == "Indexing file, please wait") {
                loading.append(" <div id =\"load\" class=\"center\">" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"textL\">Loading...</div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "<div class=\"wave\"></div>" +
                    "</div>");

            } else if (str == "Indexing complated" || currentPage >= 2) {

                document.getElementById("load").remove();
            } else { */
            //loading.append("<div class=\"textL\">Indexing complated!</div>");
            // } else
            //sysMsgAll.append("<hr>" + str + "</hr>")
        }


        //container.append(standartform);

    }
    socket.onclose = function() {
        quotation("sysMsgAll", "<p style='background-color: maroon; color:orange'>Connection Closed to WebSocket, tail stopped</p>");

    }
    socket.onerror = function(e) {
        sysMsgAll.append("<b style='color:red'>Some error occurred " + e.data.trim() + "<b>");
    }

    return socket;
}




function ParseXml(str, type) {
    var parser, xmlDoc, table, heyho;

    parser = new DOMParser();
    xmlDoc = parser.parseFromString(str, "application/xml");
    log = xmlDoc.getElementsByTagName("log");
    for (i = 0; i < log.length; i++) {

        if (type == typeMsg(log[i].getAttribute('type')) || type == "empty") {
            countRows++
            table +=
                "<tr " + "onclick=\"FuncClick(this)\"" + "  bgcolor =" + Color(log[i].getAttribute('type')) + ">" + "<td class=\"\"><span>" +
                typeMsg(log[i].getAttribute('type')) +
                "</span></td><td class=\"\"><span>" +
                log[i].getAttribute('module_name') +
                "</span></td><td class=\"\"><span>" +
                log[i].getAttribute('app_path') +
                "</span></td><td class=\"ellipsis\"><span>" +
                log[i].getAttribute('app_pid') +
                "</span></td><td class=\"\"><span>" +
                log[i].getAttribute('thread_id') +
                "</span></td><td class=\"\" data-ticks=" + timestamp(log[i].getAttribute('time')) + "><span>" +
                split_at_index(log[i].getAttribute('time')) +
                "</span></td><td class=\"\"><span>" +
                log[i].getAttribute('ulid') +
                "</span></td><td class=\"ellipsis\"><span>" +
                log[i].getAttribute('message') +
                "</span></td><td class=\"ellipsis\"><span>" +
                log[i].getAttribute('ext_message') +
                "</span></td></tr>";
        }
    }

    return table
}

function ParseXmlMap(str) {
    var parser, xmlDoc, tablemap;
    parser = new DOMParser();
    xmlDoc = parser.parseFromString(str, "text/xml");
    map = xmlDoc.getElementsByTagName("Map")[0].childNodes;

    for (i = 0; i < map.length; i++) {

        if (i % 2 != 0) {

            tablemap +=
                "<button id = " + "pagination" + " class=" + map[i].tagName + ">" + map[i].tagName + "</button>";
        }

    }
    var tablemap = tablemap.replace('undefined', '');

    return tablemap
}

function ParseCount(str) {
    var parser, xmlDoc, result
    parser = new DOMParser();
    xmlDoc = parser.parseFromString(str, "text/xml");
    content = xmlDoc.getElementsByTagName("countpage")[0];
    contentValue = content.childNodes[0];
    result = contentValue.nodeValue;
    return result
}




//23072021005653.991
//(year, monthIndex, day, hours, minutes, seconds, milliseconds) 
//23 07 2021 00 25 53.492
function timestamp(str) {
    date = str.substring(0, 2);
    month = str.substring(2, 4);
    year = str.substring(4, 8);
    hours = str.substring(8, 10);
    minutes = str.substring(10, 12)
    seconds = str.substring(12, 14);
    milliseconds = str.substring(14, 17);
    datum = new Date(Date.UTC(year, month - 1, date, hours, minutes, seconds, milliseconds));
    return datum.getTime()
}

function split_at_index(value) {
    norm = value.substring(0, 2) + "." + value.substring(2); //date
    norm = norm.substring(0, 5) + "." + norm.substring(5); //month 
    norm = norm.substring(0, 10) + " " + norm.substring(10); //year
    norm = norm.substring(0, 13) + ":" + norm.substring(13); //hours
    norm = norm.substring(0, 16) + ":" + norm.substring(16); //minutes
    norm = norm.substring(0, 19) + ":" + norm.substring(19); //seconds
    norm = norm.substring(0, 23) + " UTC" + norm.substring(23); //milliseconds
    return norm
}

function typeMsg(type) {
    if (type == "0") {
        msg = "INFO";
    } else if (type == "1") {
        msg = "DEBUG";
    } else if (type == "2") {
        msg = "WARNING";
    } else if (type == "3") {
        msg = "ERROR";
    } else if (type == "4") {
        msg = "FATAL";
    }
    return msg
}

function Color(type) {
    if (type == "0") {
        countInf++
        color = "#b4fcb5";
    } else if (type == "1") {
        countDbg++
        color = "#a0a0a0";
    } else if (type == "2") {
        countWar++
        color = "#fffc9b";
    } else if (type == "3") {
        countErr++
        color = "#fdb1b1";
    } else if (type == "4") {
        countFtl++
        color = "#b2ffb2";
    }
    countAll++
    return color
}

//PAGINATION



function PaginationButton(totalPages, maxPagesVisible = 10, currentPage = 1) {
    let pages = pageNumbers(totalPages, maxPagesVisible, currentPage);
    let currentPageBtn = null;
    const buttons = new Map();
    const disabled = {
        start: () => pages[0] === 1,
        prev: () => currentPage === 1 || currentPage > totalPages,
        end: () => pages.slice(-1)[0] === totalPages,
        next: () => currentPage >= totalPages
    }

    const frag = document.createDocumentFragment();
    const paginationButtonContainer =
        document.getElementById('test');
    paginationButtonContainer.className = 'pagtest';

    const createAndSetupButton = (label = '', cls = '', disabled = false, handleClick) => {
        const buttonElement = document.createElement('button');
        buttonElement.textContent = label;
        buttonElement.className = `page-btn ${cls}`;
        buttonElement.id = `page`
        buttonElement.disabled = disabled;
        buttonElement.addEventListener('click', e => {
            handleClick(e);
            this.update();
            paginationButtonContainer.value = currentPage;
            paginationButtonContainer.dispatchEvent(new CustomEvent('change', { detail: { currentPageBtn } }));
        });

        return buttonElement;
    }

    const onPageButtonClick = e => currentPage = Number(e.currentTarget.textContent);

    const onPageButtonUpdate = index => (btn) => {
        btn.textContent = pages[index];

        if (pages[index] === currentPage) {
            currentPageBtn.classList.remove('active');
            btn.classList.add('active');
            currentPageBtn = btn;
            currentPageBtn.focus();
            console.log("Page selected ");
            //управление выбранной кнокой (для теста кнопки)-> window.location.reload();
            setTimeout(
                () => {
                    //initWS(lastItem, statusS);
                },
                1 * 200
            );

        }
    };

    buttons.set(
        createAndSetupButton('start', 'start-page', disabled.start(), () => currentPage = 1),
        (btn) => btn.disabled = disabled.start()
    )

    buttons.set(
        createAndSetupButton('prev', 'prev-page', disabled.prev(), () => currentPage -= 1),
        (btn) => btn.disabled = disabled.prev()
    )

    pages.map((pageNumber, index) => {
        const isCurrentPage = currentPage === pageNumber;
        const button = createAndSetupButton(
            pageNumber, isCurrentPage ? 'active' : '', false, onPageButtonClick
        );

        if (isCurrentPage) {
            currentPageBtn = button;
        }

        buttons.set(button, onPageButtonUpdate(index));


    });

    buttons.set(
        createAndSetupButton('next', 'next-page', disabled.next(), () => currentPage += 1),
        (btn) => btn.disabled = disabled.next()
    )

    buttons.set(
        createAndSetupButton('end', 'end-page', disabled.end(), () => currentPage = totalPages),
        (btn) => btn.disabled = disabled.end()
    )

    buttons.forEach((_, btn) => frag.appendChild(btn));
    paginationButtonContainer.appendChild(frag);

    this.render = (container = document.body) => {
        container.appendChild(paginationButtonContainer);
    }

    this.update = (newPageNumber = currentPage) => {
        currentPage = newPageNumber;
        pages = pageNumbers(totalPages, maxPagesVisible, currentPage);
        buttons.forEach((updateButton, btn) => updateButton(btn));
    }

    this.onChange = (handler) => {
        paginationButtonContainer.addEventListener('change', handler);
    }
}

function clearBox(elementID) {
    document.getElementById(elementID).innerHTML = "";
}