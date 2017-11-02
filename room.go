package main

type room struct {
	//this channel holds all the messages
	forward chan []byte
	//this is the channel of clients wishing to join
	join chan *client
	//this is a channel of client wishing to leave
	leave chan *client
	//hold all active clients
	clients map[*client]bool
}
func (r *room) Run(){
	for{
		select {
			case client := <-r.join: r.clients[client]=true
			case client :=<- r.leave:
				delete(r.clients, client)
				close(client.send)
			case msg := <- r.forward:
				for client := range r.clients{
					client.send <- msg
				}
		}
	}
}