//
//  NetworkMonitor.swift
//  TelltaleHeart
//
//  Created by Thomas Haszard on 22/9/2024.
//

import Network


class NetworkMonitor {
    private var monitor: NWPathMonitor?
    private let queue = DispatchQueue.global(qos: .background)

    init() {
        monitor = NWPathMonitor()
        monitor?.start(queue: queue)
    }

    func isNetworkAvailable(completion: @escaping (Bool) -> Void) {
        monitor?.pathUpdateHandler = { path in
            completion(path.status == .satisfied)
        }
    }

    deinit {
        monitor?.cancel()
    }
}
