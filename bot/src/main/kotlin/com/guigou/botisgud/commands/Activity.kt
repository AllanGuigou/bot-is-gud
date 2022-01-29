package com.guigou.botisgud.commands

import com.guigou.botisgud.extensions.logger
import dev.kord.core.Kord
import dev.kord.core.event.user.VoiceStateUpdateEvent
import dev.kord.core.on
import kotlinx.coroutines.CoroutineScope

class Activity : Command {
    companion object {
        var logger = logger()
    }

    override suspend fun register(client: Kord, scope: CoroutineScope) {
        client.on<VoiceStateUpdateEvent> {
            when (state.channelId) {
                null -> logger.info("${state.userId} leave")
                else -> logger.info("${state.userId} join")
            }
        }
    }
}