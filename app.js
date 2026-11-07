var MAX_DEPTH = 3; /* parent, child, grandchild */
var items = [];
var collapsed = {};
var adding = null; /* id of the item getting a child, null for the top level */
var tree = document.getElementById("tree");
var msg = document.getElementById("msg");

function el(tag, cls, text) {
  var e = document.createElement(tag);
  e.className = cls;
  e.textContent = text;
  return e;
}

function button(cls, text, label) {
  var b = el("button", cls, text);
  b.type = "button";
  b.setAttribute("aria-label", label);
  return b;
}

function fail(err) {
  msg.textContent = err.message;
}

function api(method, body) {
  return fetch("/items", {
    method: method,
    headers: { "Content-Type": "application/json" },
    body: body && JSON.stringify(body)
  }).then(function (res) {
    if (res.ok) return res.json();
    return res.text().then(function (text) {
      throw new Error(text || "HTTP " + res.status);
    });
  });
}

/* The form for the current target is always the first input in the tree. */
function edit(id) {
  adding = id;
  render();
  tree.getElementsByTagName("input")[0].focus();
}

function form(parent_id) {
  var li = el("li", "n");
  var f = el("form", "line");
  var input = el("input", "");
  var label = parent_id === null ? "Add a task" : "Add a sub-task";

  input.placeholder = label;
  input.setAttribute("aria-label", label);
  input.onkeydown = function (e) {
    if (e.key === "Escape") edit(null);
  };
  f.onsubmit = function (e) {
    var description = input.value.trim();
    e.preventDefault();
    if (!description) return;
    api("POST", { parent_id: parent_id, description: description }).then(function (item) {
      items.push(item);
      msg.textContent = "";
      edit(adding);
    }).catch(fail);
  };

  f.appendChild(input);
  li.appendChild(el("i", "dot"));
  li.appendChild(f);
  return li;
}

function node(list, item, kids, depth) {
  var li = el("li", "n");
  var line = el("div", "line");
  var plus = button("plus", "+", "Add a sub-task");
  var ul, chev, i;

  plus.onclick = function (e) {
    e.stopPropagation();
    delete collapsed[item.id];
    edit(item.id);
  };
  line.appendChild(el("span", "txt", item.description));
  if (depth < MAX_DEPTH) line.appendChild(plus);
  li.appendChild(el("i", "dot"));
  li.appendChild(line);

  if (kids[item.id] || adding === item.id) {
    li.className = collapsed[item.id] ? "n branch closed" : "n branch";
    chev = button("chev", "›", "Expand or collapse");
    chev.setAttribute("aria-expanded", !collapsed[item.id]);
    line.appendChild(chev);
    line.onclick = function () {
      collapsed[item.id] = li.classList.toggle("closed");
      chev.setAttribute("aria-expanded", !collapsed[item.id]);
    };

    ul = el("ul", "");
    for (i = 0; kids[item.id] && i < kids[item.id].length; i++) {
      node(ul, kids[item.id][i], kids, depth + 1);
    }
    if (adding === item.id) ul.appendChild(form(item.id));
    li.appendChild(ul);
  }
  list.appendChild(li);
}

function render() {
  var kids = {};
  var i, parent;

  for (i = 0; i < items.length; i++) {
    parent = items[i].parent_id;
    if (!kids[parent]) kids[parent] = [];
    kids[parent].push(items[i]);
  }

  tree.textContent = "";
  for (i = 0; kids[null] && i < kids[null].length; i++) {
    node(tree, kids[null][i], kids, 1);
  }
  tree.appendChild(form(null));
}

edit(null);
api("GET").then(function (data) {
  items = data;
  edit(null);
}).catch(fail);
