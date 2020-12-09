package com.atlas.morg.model;

public record Monster(int worldId, int channelId, int mapId, int uniqueId, int monsterId, Integer controlCharacterId, int x, int y,
                      int fh, int stance, int team) {
}
