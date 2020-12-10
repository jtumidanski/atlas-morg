package com.atlas.morg.processor;

import java.util.Comparator;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

import com.atlas.morg.MonsterRegistry;
import com.atlas.morg.event.producer.MonsterControlEventProducer;
import com.atlas.morg.event.producer.MonsterEventProducer;
import com.atlas.morg.model.Monster;

public final class MonsterProcessor {
   private MonsterProcessor() {
   }

   public static Monster createMonster(int worldId, int channelId, int mapId, int monsterId, int x, int y, int fh, int stance,
                                       int team) {
      Monster monster = MonsterRegistry.getInstance().createMonster(worldId, channelId, mapId, monsterId, x, y, fh, stance, team);

      getControllerCandidate(worldId, channelId, mapId)
            .ifPresent(controllerId -> startControl(monster.uniqueId(), controllerId));

      MonsterEventProducer.getInstance().sendCreated(worldId, channelId, mapId, monster.uniqueId());
      return monster;
   }

   protected static Optional<Integer> getControllerCandidate(int worldId, int channelId, int mapId) {
      Map<Integer, Integer> controllerCounts = MapProcessor.getCharacterIdsInMap(worldId, channelId, mapId).stream()
            .collect(Collectors.toMap(id -> id, id -> 0));
      controllerCounts.putAll(
            MonsterRegistry.getInstance().getMonstersInMap(worldId, channelId, mapId).stream()
                  .filter(monster -> monster.controlCharacterId() != null)
                  .collect(Collectors.groupingBy(Monster::controlCharacterId, Collectors.summingInt(monster -> 1))));

      return controllerCounts.entrySet().stream()
            .min(Comparator.comparingInt(Map.Entry::getValue))
            .map(Map.Entry::getKey);
   }

   public static void findNextController(int worldId, int channelId, int mapId, int uniqueId) {
      getControllerCandidate(worldId, channelId, mapId).ifPresent(id -> startControl(uniqueId, id));
   }

   public static void startControl(int uniqueId, int characterId) {
      MonsterRegistry.getInstance().getMonster(uniqueId)
            .ifPresent(monster -> {
               if (monster.controlCharacterId() != null) {
                  stopControl(monster.controlCharacterId());
               }
               MonsterRegistry.getInstance().controlMonster(uniqueId, characterId);
               MonsterControlEventProducer.getInstance().sendControl(monster.worldId(), monster.channelId(), characterId, uniqueId);
            });
   }

   public static void stopControl(int uniqueId) {
      MonsterRegistry.getInstance().getMonster(uniqueId)
            .ifPresent(monster -> {
               MonsterRegistry.getInstance().clearControl(uniqueId);
               MonsterControlEventProducer.getInstance().clearControl(monster.worldId(), monster.channelId(),
                     monster.controlCharacterId(), uniqueId);
            });
   }
}
