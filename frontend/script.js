document.addEventListener('DOMContentLoaded', () => {
    const chatNameInput = document.getElementById('chatName');
    const chatCardsContainer = document.getElementById('chatCardsContainer');
    const chatDialog = document.getElementById('chatDialog');
    const messagesContainer = document.getElementById('messagesContainer');
    const prevPageButton = document.getElementById('prevPage');
    const nextPageButton = document.getElementById('nextPage');

    let currentPage = 1;
    let currentChatId = null;

    const fetchChats = async (name = '') => {
        try {
            const response = await fetch(`http://0.0.0.0:8080/chat/?page=1&name=${name}`, {
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
            const response = await fetch(`http://0.0.0.0:8080/chat/messages?chat_id=${chatId}&page=${page}`, {
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
        nextPageButton.disabled = messagesCount < 5; // Assuming 10 messages per page
    };

    prevPageButton.addEventListener('click', () => {
        if (currentPage > 1) {
            fetchMessages(currentChatId, currentPage - 1);
        }
    });

    nextPageButton.addEventListener('click', () => {
        fetchMessages(currentChatId, currentPage + 1);
    });

    chatNameInput.addEventListener('input', () => {
        const name = chatNameInput.value;
        fetchChats(name);
    });

    // Initial fetch
    fetchChats();
});
