classes = [];
prereqs = {};

$(document).ready(function() {
  $("#classesForm").submit(function(ev) {
    value = formatClass($("#newitem").val());

    if (value == "") {
      ev.preventDefault();
      return;
    }

    for (var i = 0; i < classes.length; i++) {
      if (classes[i] == value) {
        ev.preventDefault();
        return;
      }
    }

    classes.push(value);

    prereqs[value] = "";

    GetPrereqs(value);

    updateList();

    ev.preventDefault();
  });

  function GetPrereqs(cls) {
    $.getJSON("/api/" + cls, {}, function(data) {
      prereqs[value] = data["Prereqs"]
      updateList();
    });
  }

  function formatClass(value) {
    format = /([A-Za-z]{3}) ?([0-9]{4}[A-Z]?)/g;

    matches = format.exec(value);

    if ( matches.length == 3 ) {
      return matches[1].toUpperCase() + matches[2];
    } else if ( matches.length == 4 ) {
      return matches[1].toUpperCase() + matches[2] + matches[3];
    }

    return "";
  }

  function updateList() {
    tableObj = $("#classList");

    html = "<tr><th>Class</th><th>Prerequisites</th></tr>";

    for (var i = 0; i < classes.length; i++) {
      html += "<tr><td><span class='class'>" + classes[i] + "</span></td><td class='prereq'>" + prereqs[classes[i]] + "</td></tr>";
    }

    tableObj.html(html);
  }

  $("#classList").click(function(ev) {
    target = ev.target;

    val = $(target).text();

    newClasses = [];

    for (var i = 0; i < classes.length; i++) {
      if (classes[i] != val) {
        newClasses.push(val);
      }
    }

    classes = newClasses;

    updateList();
  })
});
