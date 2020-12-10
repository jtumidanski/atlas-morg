package com.atlas.morg.event.consumer;

import com.atlas.kafka.consumer.SimpleEventHandler;
import com.atlas.morg.MonsterRegistry;
import com.atlas.morg.processor.MonsterProcessor;
import com.atlas.mrg.constant.EventConstants;
import com.atlas.mrg.event.MapCharacterEvent;

public class MapCharacterConsumer implements SimpleEventHandler<MapCharacterEvent> {
   @Override
   public void handle(Long key, MapCharacterEvent event) {
      switch (event.type()) {
         case ENTER -> gainControl(event);
         case EXIT -> removeControl(event);
      }
   }

   private void gainControl(MapCharacterEvent event) {
      MonsterRegistry.getInstance().getMonstersInMap(event.worldId(), event.channelId(), event.mapId()).stream()
            .filter(monster -> monster.controlCharacterId() == null)
            .forEach(monster -> MonsterProcessor
                  .findNextController(event.worldId(), event.channelId(), event.mapId(), monster.uniqueId()));
   }

   private void removeControl(MapCharacterEvent event) {
      MonsterRegistry.getInstance().getMonstersInMap(event.worldId(), event.channelId(), event.mapId()).stream()
            .filter(monster -> monster.controlCharacterId() != null)
            .filter(monster -> monster.controlCharacterId() == event.characterId())
            .forEach(monster -> {
               MonsterProcessor.stopControl(monster.uniqueId());
               MonsterProcessor.findNextController(event.worldId(), event.channelId(), event.mapId(), monster.uniqueId());
            });
   }

   @Override
   public Class<MapCharacterEvent> getEventClass() {
      return MapCharacterEvent.class;
   }

   @Override
   public String getConsumerId() {
      return "Monster Registry";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return System.getenv(EventConstants.TOPIC_MAP_CHARACTER_EVENT);
   }
}
