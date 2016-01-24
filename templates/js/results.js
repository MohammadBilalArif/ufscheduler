$(document).ready(function() {

  function getData(met, unmet) {
    $.getJSON("/api/calc?" + (new Date()).getTime(), {
      "met": met,
      "unmet": unmet
    },
    function(data) {
      console.log(data);
      writeData(data);
    });
  }

  classes = [];

  function displayData() {
    html = "";

    for (var i = 0; i < classes.length; i++) {
      html += "<table class='table-hover class-list-table table'>";
      for (var cls = 0; cls < classes[i].length; cls++) {
        obj = classes[i][cls];

        html += "<tr class='classObj'><td>" + obj["Course"] + "</td><td>" + obj["Title"] + "</td><td>" + obj["Credits"] + "</td></tr>\n";
      }
      html += "</table>";
    }

    dataspot.html(html);
  }

  function writeData(data) {

    dataspot = $("#dataspot");
    html = "";

    classes = [];

    for (var i = 0; i < data.Groups.length; i++) {
      html += "<table class='table-hover class-list-table table'>";
      group = [];
      for (var cls = 0; cls < data.Groups[i].Classes.length; cls++) {
        obj = data.Groups[i].Classes[cls];

        html += "<tr class='classObj'><td class='course'>" + obj["Course"] + "</td><td class='title'>" + obj["Title"] + "</td><td class='credits'>" + obj["Credits"] + "</td></tr>\n";

        group.push({Course: obj["Course"], Title: obj["Title"], Credits: obj["Credits"]});
      }
      classes.push(group);
      html += "</table>";
    }

    dataspot.html(html);

    $("#data").css("visibility", "visible");
    $("#loading").css("visibility", "hidden");
    $("#loading").css("display","none");
  }

  met = $("#met").val()
  unmet = $("#unmet").val()

  getData(met, unmet);

  $(".course").click(function(ev) {
    alert("HERE");
  });

  $(".title").click(function(ev) {
    alert("THERE");
  });

  $(".credits").click(function(ev) {
    alert("GONE");
  });

});
