package attributes

type MonsterDataDataContainer struct {
   Data MonsterDataData `json:"data"`
}

type MonsterDataData struct {
   Id         string                `json:"id"`
   Type       string                `json:"type"`
   Attributes MonsterDataAttributes `json:"attributes"`
}

type MonsterDataAttributes struct {
   Name               string `json:"name"`
   Hp                 int    `json:"hp"`
   Mp                 int    `json:"mp"`
   Experience         int    `json:"experience"`
   Level              int    `json:"level"`
   PaDamage           int    `json:"paDamage"`
   PdDamage           int    `json:"pdDamage"`
   MaDamage           int    `json:"maDamage"`
   MdDamage           int    `json:"mdDamage"`
   Friendly           bool   `json:"friendly"`
   RemoveAfter        int    `json:"removeAfter"`
   Boss               bool   `json:"boss"`
   ExplosiveReward    bool   `json:"explosiveReward"`
   FFALoot            bool   `json:"ffaLoot"`
   Undead             bool   `json:"undead"`
   BuffToGive         int    `json:"buffToGive"`
   CarnivalPoint      int    `json:"carnivalPoint"`
   RemoveOnMiss       bool   `json:"removeOnMiss"`
   Changeable         bool   `json:"changeable"`
   TagColor           byte   `json:"tagColor"`
   TagBackgroundColor byte   `json:"tagBackgroundColor"`
   FixedStance        int    `json:"fixedStance"`
   FirstAttack        bool   `json:"firstAttack"`
   DropPeriod         int    `json:"dropPeriod"`
}
