const API = "/tasks";

let tasks = [];

fetchTasks();

async function fetchTasks() {
  try {
    const res = await fetch(API);
    if (!res.ok) throw new Error(await res.text());
    tasks = await res.json();
    render();
  } catch (e) {
    showToast("Could not load Tasks: " + e.message);
  }
}

//DOM manip :
document.getElementById("addForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const input = document.getElementById("titleInput");
  const title = input.value.trim();
  if (!title) return;
  const btn = document.getElementById("addBtn");
  btn.disabled = true;

  try {
    const res = await fetch(API, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title }),
    });
    if (!res.ok) {
      const err = await res.json();
      throw new Error(err.error || "Unknown error");
    }
    const task = await res.json();
    tasks.push(task);
    input.value = "";
    render();
    input.focus();
  } catch (e) {
    showToast("Add failed: " + e.message);
  } finally {
    btn.disabled = false;
  }
});

async function toggleTask(id) {
  try {
    const res = await fetch(`${API}/${id}`, { method: "PATCH" });
    if (!res.ok) throw new Error("Toggle Failed");
    const updated = await res.json();
    tasks = task.map((t) => (t.id === id ? updated : t));
    render();
  } catch (e) {
    showToast(e.message);
  }
}

async function deleteTask(id) {
  try {
    const res = await fetch(`${API}/${id}$`, { method: "DELETE" });
    if (!res.ok) throw new Error("Delete failed");
    tasks = task.filter((t) => t.id !== id);
    render();
  } catch (e) {
    showToast(e.message);
  }
}

function render() {
  const list = document.getElementById("taskList");
  const stats = document.getElementById("stats");

  const total = tasks.length;
  const done = tasks.filter((t) => t.done).length;

  stats.innerHTML =
    total === 0
      ? ""
      : `<span><strong>${total}</strong> task${total !== 1 ? "s" : ""}</span>` +
        `<span><strong>${done}</strong> done</span>` +
        `<span><strong>${total - done}</strong> Pending</span>`;

  if (total === 0) {
    list.innerHTML = ` <div class="empty">
          <div class="icon">📋</div>
          <p>No tasks yet — add one above.</p>
        </div>`;
    return;
  }

  list.innerHTML = tasks
    .map(
      (t) => `
      <div class="task-card${t.done ? " done" : ""}" data-id="${t.id}">
        <button class="check-btn${t.done ? " checked" : ""}"
                onclick="toggleTask('${t.id}')"
                title="${t.done ? "Mark undone" : "Mark done"}">
        </button>
        <span class="task-title">${escHtml(t.title)}</span>
        <span class="task-meta">${fmtDate(t.created_at)}</span>
        <button class="del-btn" onclick="deleteTask('${t.id}')" title="Delete">✕</button>
      </div>
    `,
    )
    .join("");
}

function escHtml(s) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

function fmtDate(iso) {
  const d = new Date(iso);
  const now = new Date();
  const diff = (now - d) / 1000;
  if (diff < 60) return "just now";
  if (diff < 3600) return Math.floor(diff / 60) + "m ago";
  if (diff < 86400) return Math.floor(diff / 3600) + "h ago";
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric" });
}

let toastTimer;

function showToast(msg) {
  const el = document.getElementById("toast");
  el.textContent = msg;
  el.classList.add("show");
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => el.classList.remove("show"), 3500);
}
