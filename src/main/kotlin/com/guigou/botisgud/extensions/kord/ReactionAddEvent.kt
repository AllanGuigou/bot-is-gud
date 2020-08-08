package com.guigou.botisgud.extensions.kord

import com.gitlab.kordlib.core.event.message.ReactionAddEvent
import io.ktor.http.Url

val ReactionAddEvent.link: Url
    get() = Url("https://discordapp.com/channels/${guildId?.value ?: "@me"}/${channelId.value}/${messageId.value}")
