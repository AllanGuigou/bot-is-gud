package com.guigou.botisgud.commands

import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.services.grpc.PresenceServiceClient
import dev.kord.core.Kord
import dev.kord.core.event.user.VoiceStateUpdateEvent
import dev.kord.core.on
import kotlinx.coroutines.CoroutineScope
import java.time.Instant

class Presence(private val presenceServiceClient: PresenceServiceClient) : Command {
    companion object {
        var logger = logger()
    }

    override suspend fun register(client: Kord, scope: CoroutineScope) {
        client.on<VoiceStateUpdateEvent> {
            val now = Instant.now()

            if (old?.channelId != null && state.channelId != null) {
                return@on
            }

            val status = when (state.channelId) {
                null -> "inactive"
                else -> "active"
            }

            logger.info("${state.userId} $status")

            presenceServiceClient.trackEvent(state.userId.value.toString(), status, now)
        }
    }
}