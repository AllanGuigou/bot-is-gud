package com.guigou.botisgud.models

import com.gitlab.kordlib.common.entity.Snowflake
import io.ktor.http.Url

data class ReminderDto(val userId: Snowflake, val message: String, val link: Url, val trigger: ReminderTrigger)
