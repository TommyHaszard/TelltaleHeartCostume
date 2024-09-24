//
//  HealthKitManger.swift
//  TelltaleHeart
//
//  Created by Thomas Haszard on 22/9/2024.
//

import HealthKit
import Network


class HealthKitManager: ObservableObject {
    let healthStore = HKHealthStore()
    let url = URL(string: "ws://192.168.1.104:8080/heartbeat")!
    @Published var isAuthorized: Bool = false
    @Published var heartRate: Int?

    func requestAuthorization() {
        let heartRateType = HKQuantityType.quantityType(forIdentifier: .heartRate)!
        
        let dataTypes: Set = [heartRateType]
        
        healthStore.requestAuthorization(toShare: nil, read: dataTypes) { success, error in
            if success {
                self.isAuthorized = true;
                self.startHeartRateQuery()
            } else if let error = error {
                print("HealthKit authorization failed: \(error.localizedDescription)")
            }
        }
    }
    
    func startHeartRateQuery() {
        let heartRateType = HKQuantityType.quantityType(forIdentifier: .heartRate)!
        
        let predicate = HKQuery.predicateForSamples(withStart: Date(), end: nil, options: .strictStartDate)
        
        let query = HKAnchoredObjectQuery(type: heartRateType, predicate: predicate, anchor: nil, limit: HKObjectQueryNoLimit) { query, samples, _, _, _ in
            self.handleHeartRateSamples(samples)
        }
        
        query.updateHandler = { query, samples, _, _, _ in
            self.handleHeartRateSamples(samples)
        }
        
        healthStore.execute(query)
    }
    
    func handleHeartRateSamples(_ samples: [HKSample]?) {
        guard let heartRateSamples = samples as? [HKQuantitySample] else { return }
        
        for sample in heartRateSamples {
            let heartRate = Int(sample.quantity.doubleValue(for: HKUnit(from: "count/min")))
            // Send this data over WebSocket
            let message = String(heartRate)
            //print(message)
            DispatchQueue.main.async {
                self.postHTTP(message)
                self.heartRate = heartRate
            }
        }
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
                return
            }

            if let httpStatus = response as? HTTPURLResponse, httpStatus.statusCode != 200 {
                // check for http errors
                print("message recieved correctly")
                print("statusCode should be 200, but is \(httpStatus.statusCode)")
            }
        }
        task.resume()
    }
}
