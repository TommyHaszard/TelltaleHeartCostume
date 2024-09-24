//
//  ExtensionDelegate.swift
//  TelltaleHeart
//
//  Created by Thomas Haszard on 22/9/2024.
//

import WatchKit
import Foundation
import HealthKit

class ExtensionDelegate: NSObject, WKExtensionDelegate {
    let healthManager = HealthKitManager()

    func applicationDidFinishLaunching() {
        // Start WebSocket and HealthKit queries
        healthManager.requestAuthorization()
    }

    func applicationWillResignActive() {
        // Optionally disconnect WebSocket
    }
}
