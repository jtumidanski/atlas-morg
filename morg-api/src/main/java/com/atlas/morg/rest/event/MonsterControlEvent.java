package com.atlas.morg.rest.event;

public record MonsterControlEvent(int worldId, int channelId, int characterId, int uniqueId, MonsterControlEventType type) {
}
