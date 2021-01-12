package com.guigou.botisgud.commands

import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.event.channel.TypingStartEvent
import com.gitlab.kordlib.core.on
import kotlinx.coroutines.CoroutineScope
import me.tatarka.inject.annotations.Inject
import java.time.Instant
import kotlin.random.Random

@Inject
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
