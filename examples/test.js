import Amqp from 'k6/x/amqp';

const url = "amqps://localhost:443";
const session = Amqp.start({
  connection_url: url,
  username: username,
  password: password,
})
export default function () {
  
  const queueName = 'TEST.QUEUE.1::TEST.QUEUE.1';
  
  Amqp.publish({
    session: session,
    queue_name: queueName,
    body: "Ping from k6 script567"
  })

  let message = Amqp.listen({
    session: session,
    queue_name: queueName,
    auto_ack: true,
  })
  console.log("Message received: " + message);
}
