//
//  TelltaleHeartApp.swift
//  TelltaleHeart Watch App
//
//  Created by Thomas Haszard on 22/9/2024.
//

import SwiftUI

@main
struct TelltaleHeart_Watch_AppApp: App {
    @StateObject private var healthKitManager = HealthKitManager()
       
       var body: some Scene {
           WindowGroup {
               if healthKitManager.isAuthorized {
                   Text("Heart Rate: \(healthKitManager.heartRate ?? 0)").font(.largeTitle)

               } else {
                   Text("Requesting Health Data Access...")
                       .onAppear {
                           // Request authorization for HealthKit access when the app launches.
                           healthKitManager.requestAuthorization()
                       }
               }
           }
       }
}
