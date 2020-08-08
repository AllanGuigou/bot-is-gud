package com.guigou.botisgud.models

import com.gitlab.kordlib.common.entity.Snowflake
import java.time.Instant

data class Reminder(val userId: Snowflake, val message: String, val link: String, val timestamp: Instant)
