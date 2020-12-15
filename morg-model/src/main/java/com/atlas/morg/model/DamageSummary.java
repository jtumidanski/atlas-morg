package com.atlas.morg.model;

public record DamageSummary(int characterId, Monster monster, long visibleDamage, long actualDamage, boolean killed) {
}
