addEventListener("DOMContentLoaded", function () {
    fetch("/api/pubkey").
        then((responce) => responce.text()).
        then((pubkey) => { document.getElementById("pubkey").innerText = pubkey; });
});

function showDialog(msg, btn) {
    let dialog = document.getElementById("dialog");
    dialog.getElementsByTagName("div")[0].innerText = msg;
    let button = dialog.getElementsByTagName("button")[0];
    if (btn) {
        button.style.display = "";
    } else {
        button.style.display = "none";
    }
    dialog.style.display = "";
}

function closeDialog() {
    let dialog = document.getElementById("dialog");
    dialog.style.display = "none";
}

function switchPane(pane) {
    let main = document.getElementById("main");
    let election = document.getElementById("election");
    main.style.display = "none";
    election.style.display = "none";
    if (pane == "election")
        election.style.display = "";
    else
        main.style.display = "";
}

function joinElection() {
    let joinStr = document.getElementById("text-join").value;
    showDialog("Joining election...", false);
    fetch("/api/election/join" + encodeURIComponent(joinStr)).
        then((responce) => responce.text()).
        then(() => { showElection(joinStr); });
}

function showElection(invStr) {
    fetch("/api/election/info/" + invStr).
        then((responce) => responce.json()).
        then((data) => {
            document.getElementById("election-title").innerText = 'Election "' + data.title + '"';
            document.getElementById("election-description").innerText = data.desription;
            closeDialog();
            switchPane("election");
        });
}
