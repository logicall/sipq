from transport import nil 

echo "hello" 

var
    transportSocket:transport.TransportSocket 

transportSocket = transport.newTransportSocket(transport.TransportType.Udp,transport.IpClass.Ip4)

