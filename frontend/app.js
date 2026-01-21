const API_URL = 'http://localhost:8080/api/v1';

let todos = [];
let editId = null;
let token = null;

window.onload = function() {
    token = localStorage.getItem('token');
    if (token) {
        showTodoSection();
        loadTodos();
    }
};

document.getElementById('loginBtn').onclick = function() {
    login();
};

document.getElementById('registerBtn').onclick = function() {
    register();
};

document.getElementById('logoutBtn').onclick = function() {
    logout();
};

document.getElementById('passwordInput').onkeypress = function(e) {
    if (e.key === 'Enter') {
        login();
    }
};

document.getElementById('addBtn').onclick = function() {
    addTodo();
};

document.getElementById('todoInput').onkeypress = function(e) {
    if (e.key === 'Enter') {
        addTodo();
    }
};

document.getElementById('searchInput').oninput = function() {
    renderTodos();
};

function register() {
    let username = document.getElementById('usernameInput').value.trim();
    let password = document.getElementById('passwordInput').value.trim();

    if (username === '' || password === '') {
        showAuthMessage('Please enter username and password!', 'error');
        return;
    }

    fetch(API_URL + '/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username: username, password: password })
    })
    .then(function(response) {
        return response.json().then(function(data) {
            if (response.ok) {
                showAuthMessage('Registration successful! Please login.', 'success');
            } else {
                showAuthMessage(data.message || 'Registration failed!', 'error');
            }
        });
    })
    .catch(function(error) {
        showAuthMessage('Connection error!', 'error');
    });
}

function login() {
    let username = document.getElementById('usernameInput').value.trim();
    let password = document.getElementById('passwordInput').value.trim();

    if (username === '' || password === '') {
        showAuthMessage('Please enter username and password!', 'error');
        return;
    }

    fetch(API_URL + '/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username: username, password: password })
    })
    .then(function(response) {
        return response.json().then(function(data) {
            if (response.ok) {
                token = data.token;
                localStorage.setItem('token', token);
                showTodoSection();
                loadTodos();
            } else {
                showAuthMessage(data.message || 'Login failed!', 'error');
            }
        });
    })
    .catch(function(error) {
        showAuthMessage('Connection error!', 'error');
    });
}

function logout() {
    token = null;
    localStorage.removeItem('token');
    todos = [];
    showLoginSection();
}

function showLoginSection() {
    document.getElementById('loginSection').style.display = 'block';
    document.getElementById('todoSection').style.display = 'none';
    document.getElementById('usernameInput').value = '';
    document.getElementById('passwordInput').value = '';
    document.getElementById('authMessage').textContent = '';
}

function showTodoSection() {
    document.getElementById('loginSection').style.display = 'none';
    document.getElementById('todoSection').style.display = 'block';
}

function showAuthMessage(message, type) {
    let msgElement = document.getElementById('authMessage');
    msgElement.textContent = message;
    msgElement.className = type;
}

function loadTodos() {
    fetch(API_URL + '/todos', {
        method: 'GET',
        headers: {
            'Authorization': 'Bearer ' + token
        }
    })
    .then(function(response) {
        if (response.status === 401) {
            logout();
            return;
        }
        return response.json();
    })
    .then(function(data) {
        if (data) {
            todos = data;
            renderTodos();
        }
    })
    .catch(function(error) {
        console.log('Error loading todos:', error);
    });
}

function addTodo() {
    let input = document.getElementById('todoInput');
    let text = input.value.trim();

    if (text === '') {
        alert('Please enter a task!');
        return;
    }

    fetch(API_URL + '/todos', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
        },
        body: JSON.stringify({ title: text })
    })
    .then(function(response) {
        if (response.status === 401) {
            logout();
            return;
        }
        return response.json();
    })
    .then(function(data) {
        if (data && data.id) {
            input.value = '';
            loadTodos();
        }
    })
    .catch(function(error) {
        console.log('Error adding todo:', error);
    });
}

function deleteTodo(id) {
    if (confirm('Are you sure you want to delete this task?')) {
        fetch(API_URL + '/todos/' + id, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        })
        .then(function(response) {
            if (response.status === 401) {
                logout();
                return;
            }
            loadTodos();
        })
        .catch(function(error) {
            console.log('Error deleting todo:', error);
        });
    }
}

function toggleComplete(id) {
    let todo = todos.find(function(t) { return t.id === id; });
    if (!todo) return;

    fetch(API_URL + '/todos/' + id, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
        },
        body: JSON.stringify({
            title: todo.title,
            completed: !todo.completed
        })
    })
    .then(function(response) {
        if (response.status === 401) {
            logout();
            return;
        }
        loadTodos();
    })
    .catch(function(error) {
        console.log('Error updating todo:', error);
    });
}

function startEdit(id) {
    editId = id;
    renderTodos();
}

function saveEdit(id) {
    let input = document.getElementById('edit-' + id);
    let newText = input.value.trim();

    if (newText === '') {
        alert('Task cannot be empty!');
        return;
    }

    let todo = todos.find(function(t) { return t.id === id; });
    if (!todo) return;

    fetch(API_URL + '/todos/' + id, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
        },
        body: JSON.stringify({
            title: newText,
            completed: todo.completed
        })
    })
    .then(function(response) {
        if (response.status === 401) {
            logout();
            return;
        }
        editId = null;
        loadTodos();
    })
    .catch(function(error) {
        console.log('Error updating todo:', error);
    });
}

function cancelEdit() {
    editId = null;
    renderTodos();
}

function renderTodos() {
    let searchText = document.getElementById('searchInput').value.toLowerCase();
    let list = document.getElementById('todoList');
    list.innerHTML = '';

    let filteredTodos = todos.filter(function(todo) {
        return todo.title.toLowerCase().includes(searchText);
    });

    filteredTodos.forEach(function(todo) {
        let li = document.createElement('li');
        li.className = 'todo-item';

        if (todo.completed) {
            li.className += ' completed';
        }

        if (editId === todo.id) {
            li.innerHTML = '<input type="text" class="edit-input" id="edit-' + todo.id + '" value="' + todo.title + '">' +
                          '<div class="todo-buttons">' +
                          '<button class="btn save-btn" onclick="saveEdit(\'' + todo.id + '\')">Save</button>' +
                          '<button class="btn cancel-btn" onclick="cancelEdit()">Cancel</button>' +
                          '</div>';
        } else {
            li.innerHTML = '<span class="todo-text' + (todo.completed ? ' completed' : '') + '">' + todo.title + '</span>' +
                          '<div class="todo-buttons">' +
                          '<button class="btn complete-btn" onclick="toggleComplete(\'' + todo.id + '\')">' + (todo.completed ? 'Undo' : 'Complete') + '</button>' +
                          '<button class="btn edit-btn" onclick="startEdit(\'' + todo.id + '\')">Edit</button>' +
                          '<button class="btn delete-btn" onclick="deleteTodo(\'' + todo.id + '\')">Delete</button>' +
                          '</div>';
        }

        list.appendChild(li);
    });
}
