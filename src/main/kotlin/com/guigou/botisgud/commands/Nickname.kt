package com.guigou.botisgud.commands

import com.gitlab.kordlib.common.entity.Snowflake
import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.behavior.edit
import com.guigou.botisgud.models.AbsoluteReminderTrigger
import com.guigou.botisgud.services.WordService
import com.guigou.botisgud.services.WordServiceImpl
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import java.time.Instant.now

data class NicknameOptions(val commandTriggerExpression: String = "* 22 * * *")

class Nickname(
    private val users: Map<Snowflake, List<Snowflake>>,
    private val options: NicknameOptions = NicknameOptions(),
    private val wordService: WordService = WordServiceImpl(),
) : Command {

    override suspend fun register(client: Kord, scope: CoroutineScope) {
        val trigger = AbsoluteReminderTrigger(options.commandTriggerExpression)

        scope.launch(Dispatchers.Default) {
            while (true) {
                delay(trigger.timestamp().toEpochMilli() - now().toEpochMilli())

                for (entry in users) {
                    for (userSnowflake in entry.value) {
                        val member = client.getUser(userSnowflake)!!.asMember(entry.key)

                        try {
                            member.edit {
                                nickname = wordService.random()
                            }
                        } catch (ex: Exception) {
                            // TODO: log exception
                        }
                    }
                }
            }
        }
    }

}
