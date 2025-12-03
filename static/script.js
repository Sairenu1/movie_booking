const BASE_URL = "http://localhost:8080";


// LOGIN
function login() {
    let data = {
        username: document.getElementById("username").value,
        password: document.getElementById("password").value
    };

    fetch(BASE_URL + "/admin/login", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(data)
    })
    .then(res => {
        if (res.ok) {
            window.location.href = "/dashboard.html";
        } else {
            document.getElementById("msg").innerText = "Invalid credentials";
        }
    });
}



// DASHBOARD STATS
function loadDashboard() {
    fetch(BASE_URL + "/admin/dashboard/stats")
    .then(res => res.json())
    .then(data => {
        document.getElementById("total").innerText = data.total_bookings;
        document.getElementById("active").innerText = data.active_bookings;
        document.getElementById("inactive").innerText = data.inactive_bookings;
    });
}



// MOVIES
function loadMovies() {
    fetch(BASE_URL + "/admin/movies")
    .then(res => res.json())
    .then(movies => {
        let table = document.getElementById("movieTable");

        movies.forEach(m => {
            let row = `<tr>
                <td>${m.id}</td>
                <td>${m.title}</td>
                <td>${m.movie_number}</td>
                <td>${m.genre}</td>
                <td>${m.duration}</td>
                <td>
                    <button class="btn" onclick="deleteMovie(${m.id})">Delete</button>
                </td>
            </tr>`;
            table.innerHTML += row;
        });
    });
}

function addMovie() {
    let movie = {
        title: document.getElementById("title").value,
        movie_number: document.getElementById("number").value,
        genre: document.getElementById("genre").value,
        duration: Number(document.getElementById("duration").value)
    };

    fetch(BASE_URL + "/admin/movies", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(movie)
    })
    .then(() => {
        alert("Movie added!");
        location.reload();
    });
}

function deleteMovie(id) {
    fetch(BASE_URL + "/admin/movies/" + id, { method: "DELETE" })
    .then(() => {
        alert("Movie deleted!");
        location.reload();
    });
}



// BOOKINGS
function loadBookings() {
    fetch(BASE_URL + "/admin/bookings")
    .then(res => res.json())
    .then(bookings => {
        let table = document.getElementById("bookingTable");

        bookings.forEach(b => {
            let status = b.is_active ? "Active" : "Inactive";

            let row = `<tr>
                <td>${b.id}</td>
                <td>${b.movie}</td>
                <td>${b.user}</td>
                <td>${b.seat}</td>
                <td>${status}</td>
                <td>
                    <button class="btn" onclick="deactivate(${b.id})">Deactivate</button>
                    <button class="btn" onclick="restore(${b.id})">Restore</button>
                </td>
            </tr>`;
            table.innerHTML += row;
        });
    });
}

function deactivate(id) {
    fetch(BASE_URL + `/admin/bookings/${id}/deactivate`, {method: "PATCH"})
    .then(() => location.reload());
}

function restore(id) {
    fetch(BASE_URL + `/admin/bookings/${id}/restore`, {method: "PATCH"})
    .then(() => location.reload());
}
