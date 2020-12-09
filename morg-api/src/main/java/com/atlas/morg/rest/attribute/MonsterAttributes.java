package com.atlas.morg.rest.attribute;

import rest.AttributeResult;

public record MonsterAttributes(Integer mapId, Integer monsterId, Integer controlCharacterId, Integer x, Integer y, Integer fh,
                                Integer stance, Integer team) implements AttributeResult {
}
