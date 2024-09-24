//
//  WebSocketClient.swift
//  TelltaleHeart
//
//  Created by Thomas Haszard on 22/9/2024.
//

import Foundation

class HttpSender {
    var webSocketTask: URLSessionWebSocketTask?
    var urlSession: URLSession?
    
    func setup() {
        let url = URL(string: "ws://192.168.1.104:8080/heartbeat")!
    }
    
    func postHTTP(_ heartbeat: String) {
        let store = [
            "heartBeat" : heartbeat
        ]
        let jsonData = try? JSONSerialization.data(withJSONObject: store, options: JSONSerialization.WritingOptions.prettyPrinted)

        var request = URLRequest(url: URL(string: "http://192.168.1.104:8080/heartbeat")!)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpMethod = "POST"
        request.httpBody = jsonData
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            guard let data = data, error == nil else {
                // check for fundamental networking error
                print("error=\(error!)")
                return
            }

            if let httpStatus = response as? HTTPURLResponse, httpStatus.statusCode != 200 {
                // check for http errors
                print("statusCode should be 200, but is \(httpStatus.statusCode)")
                print("response = \(response!)")
            }

            do {
                if let json = try JSONSerialization.jsonObject(with: data, options: .mutableContainers) as? [String: Any] {
                    print(json)
                }
            } catch let error {
                print(error.localizedDescription)
            }
        }
        task.resume()
    }
}
