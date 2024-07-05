document.addEventListener('DOMContentLoaded', () => {
    const chatNameInput = document.getElementById('chatName');
    const chatCardsContainer = document.getElementById('chatCardsContainer');
    const chatDialog = document.getElementById('chatDialog');
    const messagesContainer = document.getElementById('messagesContainer');
    const prevPageButton = document.getElementById('prevPage');
    const nextPageButton = document.getElementById('nextPage');
    const messageInput = document.getElementById('messageInput');
    const sendMessageButton = document.getElementById('sendMessage');
    const leaveChatButton = document.getElementById('leaveChat');
    const createChatButton = document.getElementById('createChat');
    const errorContainer = document.getElementById('errorContainer');

    let currentPage = 1;
    let currentChatId = null;
    let socket = null;

    const fetchChats = async (name = '') => {
        try {
            const response = await fetch(`/chat/?page=1&name=${name}`, {
                headers: {
                    'accept': 'application/json'
                }
            });
            const data = await response.json();
            populateCards(data.chats);
        } catch (error) {
            console.error('Error fetching chats:', error);
        }
    };

    const fetchMessages = async (chatId, page = 1) => {
        try {
            const response = await fetch(`/chat/messages?chat_id=${chatId}&page=${page}`, {
                headers: {
                    'accept': 'application/json'
                }
            });
            const data = await response.json();
            populateMessages(data.messages);
            currentPage = page;
            currentChatId = chatId;
            updatePaginationButtons(data.messages.length);
        } catch (error) {
            console.error('Error fetching messages:', error);
        }
    };

    const populateCards = (chats) => {
        chatCardsContainer.innerHTML = '';
        chats.forEach(chat => {
            const card = document.createElement('div');
            card.className = 'chat-card';
            card.textContent = chat.name;
            card.dataset.chatId = chat.id;
            card.addEventListener('click', () => {
                fetchMessages(chat.id);
                chatDialog.classList.remove('hidden');
                connectWebSocket(chat.id);
            });
            chatCardsContainer.appendChild(card);
        });
    };

    const populateMessages = (messages) => {
        messagesContainer.innerHTML = '';
        messages = messages.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));
        messages.forEach(message => {
            const messageElement = document.createElement('div');
            messageElement.className = 'message';
            messageElement.textContent = `${message.timestamp}, ${message.sender}: ${message.content}`;
            messagesContainer.appendChild(messageElement);
        });
    };

    const updatePaginationButtons = (messagesCount) => {
        prevPageButton.disabled = currentPage === 1;
        nextPageButton.disabled = messagesCount < 10; // Assuming 10 messages per page
    };

    const connectWebSocket = (chatId) => {
        if (socket) {
            socket.close();
        }
        socket = new WebSocket(`ws://0.0.0.0:8080/ws/join_chat?id=${chatId}`);

        socket.onopen = () => {
            console.log('WebSocket connection established');
        };

        socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            const messageElement = document.createElement('div');
            messageElement.className = 'message';
            messageElement.textContent = `${message.timestamp}, ${message.sender}: ${message.content}`;
            messagesContainer.appendChild(messageElement);
        };

        socket.onclose = (event) => {
            if (event.wasClean) {
                console.log(`WebSocket connection closed cleanly, code=${event.code}, reason=${event.reason}`);
            } else {
                console.error('WebSocket connection closed unexpectedly');
            }
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    };

    sendMessageButton.addEventListener('click', () => {
        const message = messageInput.value;
        if (message && socket) {
            let socketMessage = JSON.stringify(message).replace("\"", "")
            socketMessage = socketMessage.split(/(?:)/u).reverse().join("").replace("\"", "").split(/(?:)/u).reverse().join("");

            socket.send(socketMessage);

            // Отображение отправленного сообщения на экране
            const messageElement = document.createElement('div');
            messageElement.className = 'message';
            messageElement.textContent = `${new Date()} Вы: ${message}`;
            messagesContainer.appendChild(messageElement);

            messageInput.value = '';
        }
    });

    leaveChatButton.addEventListener('click', () => {
        if (socket) {
            socket.close();
            socket = null;
            chatDialog.classList.add('hidden');
            messagesContainer.innerHTML = '';
        }
    });

    createChatButton.addEventListener('click', async () => {
        const chatName = chatNameInput.value;
        if (chatName) {
            try {
                const response = await fetch(`/chat/?name=${chatName}`, {
                    method: 'POST',
                    headers: {
                        'accept': 'application/json'
                    }
                });
                if (!response.ok) {
                    let payload = await response.json()
                    throw new Error('Failed to create chat' + payload['err']);
                }
                alert("Chat created successfully!")
                const data = await response.json();
                const newChatId = data.chat_id;
                fetchMessages(newChatId);
                chatDialog.classList.remove('hidden');
                connectWebSocket(newChatId);
                errorContainer.classList.add('hidden');
            } catch (error) {
                console.error('Error creating chat:', error);
                alert(error.message)
            }
        }
    });

    chatNameInput.addEventListener('input', () => {
        const name = chatNameInput.value;
        fetchChats(name);
    });

    prevPageButton.addEventListener('click', () => {
        if (currentPage > 1) {
            fetchMessages(currentChatId, currentPage - 1);
        }
    });

    nextPageButton.addEventListener('click', () => {
        fetchMessages(currentChatId, currentPage + 1);
    });

    // Initial fetch
    fetchChats();
});
