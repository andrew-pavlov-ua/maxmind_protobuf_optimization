package services

// Generate a random IPv4 address

// func randomIPv4() net.IP {
// 	return net.IPv4(
// 		byte(rand.Intn(256)),
// 		byte(rand.Intn(256)),
// 		byte(rand.Intn(256)),
// 		byte(rand.Intn(256)),
// 	)
// }

// // Generate a random IPv6 address
// func randomIPv6() net.IP {
// 	ip := make(net.IP, net.IPv6len)
// 	crand.Read(ip)
// 	return ip
// }

// // Generate a list of random IPs (IPv4 & IPv6)
// func generateIPList(count int) []net.IP {
// 	ips := make([]net.IP, count)
// 	for i := range count {
// 		if rand.Intn(2) == 0 {
// 			ips[i] = randomIPv4()
// 		} else {
// 			ips[i] = randomIPv6()
// 		}
// 	}
// 	return ips
// }
