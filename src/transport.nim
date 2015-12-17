from net import nil 
from nativesockets as raw import nil 

type 
  TransportType* {.pure.} = enum
    Tcp
    Udp
    Tls
    Sctp

  IpClass* {.pure.} = enum 
    Ip4
    Ip6 

  TransportSocket* = object
    socket:net.Socket


proc newTransportSocket* (transportType: TransportType; ipClass:IpClass): TransportSocket = 
  var 
    protocol: raw.Protocol
    domain: raw.Domain
    sockType: raw.SockType 

   
  if ipClass == IpClass.Ip4:
    domain = raw.AF_INET
  else:
    domain = raw.AF_INET6 
      
  case transportType:
    of TransportType.Udp:
      sockType = raw.SOCK_DGRAM 
      protocol = raw.IPPROTO_UDP 
    else:
      assert(false) 
    
  var socket:net.Socket

  try :
    socket = net.newSocket(domain, sockType,protocol)
  except OSError:
    raise 
    
  result.socket = socket 

