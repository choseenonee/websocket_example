document.addEventListener('DOMContentLoaded', () => {
    const chatNameInput = document.getElementById('chatName');
    const chatCardsContainer = document.getElementById('chatCardsContainer');

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

    const populateCards = (chats) => {
        chatCardsContainer.innerHTML = '';
        chats.forEach(chat => {
            const card = document.createElement('div');
            card.className = 'chat-card';
            card.textContent = chat.name;
            card.dataset.chatId = chat.id;
            chatCardsContainer.appendChild(card);
        });
    };

    chatNameInput.addEventListener('input', () => {
        const name = chatNameInput.value;
        fetchChats(name);
    });

    // Initial fetch
    fetchChats();
});
