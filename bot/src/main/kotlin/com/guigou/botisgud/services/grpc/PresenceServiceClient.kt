package com.guigou.botisgud.services.grpc

import com.botisgud.presence.*
import com.google.protobuf.Timestamp
import com.guigou.botisgud.extensions.logger
import io.grpc.ManagedChannel
import io.grpc.StatusException
import java.io.Closeable
import java.util.concurrent.TimeUnit
import kotlinx.coroutines.runBlocking
import java.time.Instant

class PresenceServiceClient(val channel: ManagedChannel) : Closeable {
    companion object {
        var logger = logger()
    }

    private val stub: PresenceServiceGrpcKt.PresenceServiceCoroutineStub =
        PresenceServiceGrpcKt.PresenceServiceCoroutineStub(channel)

    fun trackEvent(u: String, s: String, t: Instant) = runBlocking {
        val request =
            eventRequest { user = u; status = s; timestamp = Timestamp.newBuilder().setSeconds(t.epochSecond).build() }
        try {
            stub.trackEvent(request)
        } catch (e: StatusException) {
            logger.error("RPC failed: ${e.status}")
        }
    }

    override fun close() {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS)
    }
}