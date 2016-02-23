$(document).ready(function() {

  function getData(met, unmet) {
    $.post("/api/calc?" + (new Date()).getTime(), {
      "met": met,
      "unmet": unmet
    },
    function(data) {
      writeData(JSON.parse(data));
    });
  }

  classes = [];

  function displayData() {
    html = "";

    for (var i = 0; i < classes.length; i++) {
      html += "<table class='table-hover table'>";
      html += "<tr><th>Course ID</th><th>Title</th><th>Credit Hours</th></tr>"
      for (var cls = 0; cls < classes[i].length; cls++) {
        obj = classes[i][cls];

        html += "<tr class='classObj'><td>" + obj["Course"] + "</td><td>" + obj["Title"] + "</td><td>" + obj["Credits"] + "</td></tr>\n";
      }
      html += "</table>";
    }

    dataspot.html(html);
  }

  function findObject( courseID ) {
    for (var i = 0; i < classes.length; i++) {
      for (var j = 0; j < classes[i].length; j++) {
        if ( classes[i][j].Course == courseID )
          return classes[i][j];
      }
    }

    return {};
  }



  function writeData(data) {

    dataspot = $("#dataspot");
    html = "";

    classes = [];

    for (var i = 0; i < data.Groups.length; i++) {
      html += "<table class='table table-striped'>";
      html += "<tr><th>Course ID</th><th>Title</th><th>Credit Hours</th></tr>"
      group = [];
      for (var cls = 0; cls < data.Groups[i].Classes.length; cls++) {
        obj = data.Groups[i].Classes[cls];

        html += "<tr class='" + obj["Course"] + "'><td class='course'>" + obj["Course"] + "</td><td class='title'>" + obj["Title"] + "</td><td class='credits'>" + obj["Credits"] + "</td></tr>\n";

        group.push({ "Course": obj["Course"], "Title": obj["Title"], "Credits": obj["Credits"], "Prereqs": obj["Prereqs"] });
      }
      classes.push(group);
      html += "</table>";
    }

    dataspot.html(html);
/*
    for (var j = 0; j < classes.length; j++ ) {
      createTable(classes[j]);
    }
*/
    $("table").delegate("tr", "mouseenter", function() {
      var tr = $(this).closest('tr');

      $('table tr').removeClass("highlight");

      var obj = findObject( tr.children(".course").text() );

      if (obj != {}) {
        for (var i = 0; i < obj.Prereqs.length; i++) {
          prereq = obj.Prereqs[i];

          console.log("prereq: ", prereq);

          $("."+prereq).addClass("highlight");
        }
      } else {
      }
    });

    $("#data").css("visibility", "visible");
    $("#loading").css("visibility", "hidden");
    $("#loading").css("display","none");
  }

  met = $("#met").val()
  unmet = $("#unmet").val()

  getData(met, unmet);

});
