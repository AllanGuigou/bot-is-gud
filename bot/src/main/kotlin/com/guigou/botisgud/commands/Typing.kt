package com.guigou.botisgud.commands

import dev.kord.core.Kord
import dev.kord.core.event.channel.TypingStartEvent
import dev.kord.core.on
import com.guigou.botisgud.services.grpc.HelloWorldClient
import dev.kord.core.event.message.MessageCreateEvent
import io.grpc.ManagedChannelBuilder
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.asCoroutineDispatcher
import kotlinx.coroutines.asExecutor
import java.time.Instant
import java.util.concurrent.Executors
import kotlin.random.Random

class Typing : Command {
    private var triggered: Instant = Instant.MIN;

    override suspend fun register(client: Kord, scope: CoroutineScope) {
        client.on<TypingStartEvent> {
            if (user.asUser().isBot == true) {
                return@on
            }

            if (triggered.isAfter(Instant.now().minusSeconds(60))) {
                return@on
            }

            if (Random.nextInt(0, 100) < 20) {
                return@on
            }

            channel.type()
            triggered = Instant.now()
        }
    }
}
