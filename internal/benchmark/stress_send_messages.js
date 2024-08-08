import http from 'k6/http';
import { check, fail } from 'k6';
import { SharedArray } from 'k6/data';
import ws from 'k6/ws';

const ChatsAmount = 20
const ChatFirstID = 959

const chatIDs = new SharedArray('some name', function () {
  const chatIDs = [];
  for (let i = 0; i < ChatsAmount; i++) {
    chatIDs.push(ChatFirstID + i)
  }

  return chatIDs; // must be an array
});

export const options = {
  duration: '1m', vus: 300
};

export default function () {
  const chatID = chatIDs[Math.floor(__VU-1)]

  const url = `ws://95.84.137.217:3002/ws/join_chat?id=${chatID}`;

  const res = ws.connect(url, null, function (socket) {
    socket.on('open', () => socket.send(`hello world! ${Date.now()}`));
    socket.setInterval(function timeout() {
      socket.send(`hello world! ${Date.now()}`);
    }, 100);

    socket.setTimeout(function () {
      socket.close();
    }, 60000);
  });

  const err = check(res, { 'status is 101': (r) => r && r.status === 101 });
  if (!err) {
    console.error(`err on opening webSocket ${chatID}`)
  }
}
