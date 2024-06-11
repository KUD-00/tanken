import { OhVueIcon, addIcons } from "oh-vue-icons";
import {
  CoHeart,
  CoBookmark,
  CoCommentBubble,
  CoGithub,
  CoGoogle,
  CoOpenstreetmap,
  CoGlobeAlt,
  LaUserCircleSolid,
  CoWalk,
  CoFlagAlt,
  FaWalking,
  CoThumbUp,
  CoRestaurant,
  RiShoppingCartLine,
  MdBackpackOutlined,
  IoRestaurantOutline,
  CoPencil,
  MdInsertphotoOutlined,
  BiSearch,
} from "oh-vue-icons/icons";

addIcons(
  CoHeart,
  CoBookmark,
  CoCommentBubble,
  CoGithub,
  CoGoogle,
  CoOpenstreetmap,
  CoGlobeAlt,
  LaUserCircleSolid,
  CoWalk,
  CoFlagAlt,
  FaWalking,
  CoThumbUp,
  CoRestaurant,
  RiShoppingCartLine,
  MdBackpackOutlined,
  IoRestaurantOutline,
  CoPencil,
  MdInsertphotoOutlined,
  BiSearch,
);

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.component("v-icon", OhVueIcon);
});
