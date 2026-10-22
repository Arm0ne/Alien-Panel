import{$i as e,$n as t,Cr as n,Dr as r,Er as i,Fi as a,Gi as o,Hr as s,In as c,Jn as l,Ki as u,Kn as d,Li as f,Lr as p,Mi as m,Mr as h,Or as g,Pn as _,Rr as v,Sr as y,Ur as b,Vr as x,Yi as S,Zi as C,Zn as ee,_a as w,ai as te,an as T,br as E,ca as D,ci as O,cr as k,di as A,dn as ne,en as j,fi as M,fr as N,gn as P,ha as F,ia as I,jn as L,ki as R,kn as z,li as B,ln as re,mr as ie,na as V,ni as H,nr as U,oi as W,on as G,pi as K,qi as q,qn as ae,qr as J,ri as oe,sa as Y,tn as se,tr as ce,ua as le,ui as X,un as ue,wr as Z,yr as Q}from"./router-BdUGwzn9.js";import{o as de,t as fe}from"./focus-detector-CihS7aiz.js";import{a as pe,r as me,t as he}from"./light-9-WouATb.js";import{n as ge}from"./data-freshness-B-gZruil.js";var $=`v-hidden`,_e=y(`[v-hidden]`,{display:`none!important`}),ve=o({name:`Overflow`,props:{getCounter:Function,getTail:Function,updateCounter:Function,onUpdateCount:Function,onUpdateOverflow:Function},setup(t,{slots:r}){let i=F(null),a=F(null);function o(e){let{value:n}=i,{getCounter:o,getTail:s}=t,c;if(c=o===void 0?a.value:o(),!n||!c)return;c.hasAttribute($)&&c.removeAttribute($);let{children:l}=n;if(e.showAllItemsBeforeCalculate)for(let e of l)e.hasAttribute($)&&e.removeAttribute($);let u=n.offsetWidth,d=[],f=r.tail?s?.():null,p=f?f.offsetWidth:0,m=!1,h=n.children.length-+!!r.tail;for(let e=0;e<h-1;++e){if(e<0)continue;let n=l[e];if(m){n.hasAttribute($)||n.setAttribute($,``);continue}else n.hasAttribute($)&&n.removeAttribute($);let r=n.offsetWidth;if(p+=r,d[e]=r,p>u){let{updateCounter:n}=t;for(let r=e;r>=0;--r){let i=h-1-r;n===void 0?c.textContent=`${i}`:n(i);let a=c.offsetWidth;if(p-=d[r],p+a<=u||r===0){m=!0,e=r-1,f&&(e===-1?(f.style.maxWidth=`${u-a}px`,f.style.boxSizing=`border-box`):f.style.maxWidth=``);let{onUpdateCount:n}=t;n&&n(i);break}}}}let{onUpdateOverflow:g}=t;m?g!==void 0&&g(!0):(g!==void 0&&g(!1),c.setAttribute($,``))}let s=Z();return _e.mount({id:`vueuc/overflow`,head:!0,anchorMetaName:n,ssr:s}),e(()=>o({showAllItemsBeforeCalculate:!1})),{selfRef:i,counterRef:a,sync:o}},render(){let{$slots:e}=this;return S(()=>this.sync({showAllItemsBeforeCalculate:!1})),u(`div`,{class:`v-overflow`,ref:`selfRef`},[I(e,`default`),e.counter?e.counter():u(`span`,{style:{display:`inline-block`},ref:`counterRef`}),e.tail?e.tail():null])}});function ye(t,n){n&&(e(()=>{let{value:e}=t;e&&Q.registerHandler(e,n)}),Y(t,(e,t)=>{t&&Q.unregisterHandler(t)},{deep:!1}),C(()=>{let{value:e}=t;e&&Q.unregisterHandler(e)}))}function be(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var xe=o({name:`Checkmark`,render(){return u(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},u(`g`,{fill:`none`},u(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})))}}),Se=o({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:t,labelFieldRef:n,nodePropsRef:r}=q(v);return{labelField:n,nodeProps:r,renderLabel:e,renderOption:t}},render(){let{clsPrefix:e,renderLabel:t,renderOption:n,nodeProps:r,tmNode:{rawNode:i}}=this,a=r?.(i),o=t?t(i,!1):U(i[this.labelField],i,!1),s=u(`div`,Object.assign({},a,{class:[`${e}-base-select-group-header`,a?.class]}),o);return i.render?i.render({node:s,option:i}):n?n({node:s,option:i,selected:!1}):s}});function Ce(e,t){return u(R,{name:`fade-in-scale-up-transition`},{default:()=>e?u(z,{clsPrefix:t,class:`${t}-base-select-option__check`},{default:()=>u(xe)}):null})}var we=o({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:t,pendingTmNodeRef:n,multipleRef:r,valueSetRef:i,renderLabelRef:a,renderOptionRef:o,labelFieldRef:s,valueFieldRef:c,showCheckmarkRef:l,nodePropsRef:u,handleOptionClick:d,handleOptionMouseEnter:f}=q(v),p=J(()=>{let{value:t}=n;return t?e.tmNode.key===t.key:!1});function m(t){let{tmNode:n}=e;n.disabled||d(t,n)}function h(t){let{tmNode:n}=e;n.disabled||f(t,n)}function g(t){let{tmNode:n}=e,{value:r}=p;n.disabled||r||f(t,n)}return{multiple:r,isGrouped:J(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:l,nodeProps:u,isPending:p,isSelected:J(()=>{let{value:n}=t,{value:a}=r;if(n===null)return!1;let o=e.tmNode.rawNode[c.value];if(a){let{value:e}=i;return e.has(o)}else return n===o}),labelField:s,renderLabel:a,renderOption:o,handleMouseMove:g,handleMouseEnter:h,handleClick:m}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:n,isPending:r,isGrouped:i,showCheckmark:a,nodeProps:o,renderOption:s,renderLabel:c,handleClick:l,handleMouseEnter:d,handleMouseMove:f}=this,p=Ce(n,e),m=c?[c(t,n),a&&p]:[U(t[this.labelField],t,n),a&&p],h=o?.(t),g=u(`div`,Object.assign({},h,{class:[`${e}-base-select-option`,t.class,h?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:n,[`${e}-base-select-option--grouped`]:i,[`${e}-base-select-option--pending`]:r,[`${e}-base-select-option--show-checkmark`]:a}],style:[h?.style||``,t.style||``],onClick:be([l,h?.onClick]),onMouseenter:be([d,h?.onMouseenter]),onMousemove:be([f,h?.onMousemove])}),u(`div`,{class:`${e}-base-select-option__content`},m));return t.render?t.render({node:g,option:t,selected:n}):s?s({node:g,option:t,selected:n}):g}}),Te=B(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[B(`scrollbar`,`
 max-height: var(--n-height);
 `),B(`virtual-list`,`
 max-height: var(--n-height);
 `),B(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[X(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),B(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),B(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),X(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),X(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),X(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),X(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),B(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),B(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[A(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),O(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),O(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),A(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),A(`pending`,[O(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),A(`selected`,`
 color: var(--n-option-text-color-active);
 `,[O(`&::before`,`
 background-color: var(--n-option-color-active);
 `),A(`pending`,[O(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),A(`disabled`,`
 cursor: not-allowed;
 `,[M(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),A(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),X(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[T({enterScale:`0.5`})])])]),Ee=o({name:`InternalSelectMenu`,props:Object.assign(Object.assign({},L.props),{clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function}),setup(t){let{mergedClsPrefixRef:n,mergedRtlRef:r,mergedComponentPropsRef:i}=l(t),a=_(`InternalSelectMenu`,r,n),o=L(`InternalSelectMenu`,`-internal-select-menu`,Te,pe,t,w(t,`clsPrefix`)),s=F(null),c=F(null),u=F(null),d=f(()=>t.treeMate.getFlattenedNodes()),m=f(()=>ue(d.value)),h=F(null);function g(){let{treeMate:e}=t,n=null,{value:r}=t;r===null?n=e.getFirstAvailableNode():(n=t.multiple?e.getNode((r||[])[(r||[]).length-1]):e.getNode(r),(!n||n.disabled)&&(n=e.getFirstAvailableNode())),B(n||null)}function y(){let{value:e}=h;e&&!t.treeMate.getNode(e.key)&&(h.value=null)}let b;Y(()=>t.show,e=>{e?b=Y(()=>t.treeMate,()=>{t.resetMenuOnOptionsChange?(t.autoPending?g():y(),S(re)):y()},{immediate:!0}):b?.()},{immediate:!0}),C(()=>{b?.()});let x=f(()=>H(o.value.self[K(`optionHeight`,t.size)])),ee=f(()=>oe(o.value.self[K(`padding`,t.size)])),te=f(()=>t.multiple&&Array.isArray(t.value)?new Set(t.value):new Set),T=f(()=>{let e=d.value;return e&&e.length===0}),E=f(()=>i?.value?.Select?.renderEmpty);function D(e){let{onToggle:n}=t;n&&n(e)}function O(e){let{onScroll:n}=t;n&&n(e)}function k(e){var t;(t=u.value)==null||t.sync(),O(e)}function A(){var e;(e=u.value)==null||e.sync()}function ne(){let{value:e}=h;return e||null}function j(e,t){t.disabled||B(t,!1)}function M(e,t){t.disabled||D(t)}function N(e){var n;W(e,`action`)||(n=t.onKeyup)==null||n.call(t,e)}function P(e){var n;W(e,`action`)||(n=t.onKeydown)==null||n.call(t,e)}function I(e){var n;(n=t.onMousedown)==null||n.call(t,e),!t.focusable&&e.preventDefault()}function R(){let{value:e}=h;e&&B(e.getNext({loop:!0}),!0)}function z(){let{value:e}=h;e&&B(e.getPrev({loop:!0}),!0)}function B(e,t=!1){h.value=e,t&&re()}function re(){var e,n;let r=h.value;if(!r)return;let i=m.value(r.key);i!==null&&(t.virtualScroll?(e=c.value)==null||e.scrollTo({index:i}):(n=u.value)==null||n.scrollTo({index:i,elSize:x.value}))}function ie(e){var n;s.value?.contains(e.target)&&((n=t.onFocus)==null||n.call(t,e))}function U(e){var n;s.value?.contains(e.relatedTarget)||(n=t.onBlur)==null||n.call(t,e)}V(v,{handleOptionMouseEnter:j,handleOptionClick:M,valueSetRef:te,pendingTmNodeRef:h,nodePropsRef:w(t,`nodeProps`),showCheckmarkRef:w(t,`showCheckmark`),multipleRef:w(t,`multiple`),valueRef:w(t,`value`),renderLabelRef:w(t,`renderLabel`),renderOptionRef:w(t,`renderOption`),labelFieldRef:w(t,`labelField`),valueFieldRef:w(t,`valueField`)}),V(p,s),e(()=>{let{value:e}=u;e&&e.sync()});let G=f(()=>{let{size:e}=t,{common:{cubicBezierEaseInOut:n},self:{height:r,borderRadius:i,color:a,groupHeaderTextColor:s,actionDividerColor:c,optionTextColorPressed:l,optionTextColor:u,optionTextColorDisabled:d,optionTextColorActive:f,optionOpacityDisabled:p,optionCheckColor:m,actionTextColor:h,optionColorPending:g,optionColorActive:_,loadingColor:v,loadingSize:y,optionColorActivePending:b,[K(`optionFontSize`,e)]:x,[K(`optionHeight`,e)]:S,[K(`optionPadding`,e)]:C}}=o.value;return{"--n-height":r,"--n-action-divider-color":c,"--n-action-text-color":h,"--n-bezier":n,"--n-border-radius":i,"--n-color":a,"--n-option-font-size":x,"--n-group-header-text-color":s,"--n-option-check-color":m,"--n-option-color-pending":g,"--n-option-color-active":_,"--n-option-color-active-pending":b,"--n-option-height":S,"--n-option-opacity-disabled":p,"--n-option-text-color":u,"--n-option-text-color-active":f,"--n-option-text-color-disabled":d,"--n-option-text-color-pressed":l,"--n-option-padding":C,"--n-option-padding-left":oe(C,`left`),"--n-option-padding-right":oe(C,`right`),"--n-loading-color":v,"--n-loading-size":y}}),{inlineThemeDisabled:q}=t,J=q?ae(`internal-select-menu`,f(()=>t.size[0]),G,t):void 0,se={selfRef:s,next:R,prev:z,getPendingTmNode:ne};return ye(s,t.onResize),Object.assign({mergedTheme:o,mergedClsPrefix:n,rtlEnabled:a,virtualListRef:c,scrollbarRef:u,itemSize:x,padding:ee,flattenedNodes:d,empty:T,mergedRenderEmpty:E,virtualListContainer(){let{value:e}=c;return e?.listElRef},virtualListContent(){let{value:e}=c;return e?.itemsElRef},doScroll:O,handleFocusin:ie,handleFocusout:U,handleKeyUp:N,handleKeyDown:P,handleMouseDown:I,handleVirtualListResize:A,handleVirtualListScroll:k,cssVars:q?void 0:G,themeClass:J?.themeClass,onRender:J?.onRender},se)},render(){let{$slots:e,virtualScroll:n,clsPrefix:r,mergedTheme:i,themeClass:a,onRender:o}=this;return o?.(),u(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:[`${r}-base-select-menu`,`${r}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${r}-base-select-menu--rtl`,a,this.multiple&&`${r}-base-select-menu--multiple`],style:this.cssVars,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},ce(e.header,e=>e&&u(`div`,{class:`${r}-base-select-menu__header`,"data-header":!0,key:`header`},e)),this.loading?u(`div`,{class:`${r}-base-select-menu__loading`},u(P,{clsPrefix:r,strokeWidth:20})):this.empty?u(`div`,{class:`${r}-base-select-menu__empty`,"data-empty":!0},t(e.empty,()=>[this.mergedRenderEmpty?.call(this)||u(G,{theme:i.peers.Empty,themeOverrides:i.peerOverrides.Empty,size:this.size})])):u(ne,Object.assign({ref:`scrollbarRef`,theme:i.peers.Scrollbar,themeOverrides:i.peerOverrides.Scrollbar,scrollable:this.scrollable,container:n?this.virtualListContainer:void 0,content:n?this.virtualListContent:void 0,onScroll:n?void 0:this.doScroll},this.scrollbarProps),{default:()=>n?u(de,{ref:`virtualListRef`,class:`${r}-virtual-list`,items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?u(Se,{key:e.key,clsPrefix:r,tmNode:e}):e.ignored?null:u(we,{clsPrefix:r,key:e.key,tmNode:e})}):u(`div`,{class:`${r}-base-select-menu-option-wrapper`,style:{paddingTop:this.padding.top,paddingBottom:this.padding.bottom}},this.flattenedNodes.map(e=>e.isGroup?u(Se,{key:e.key,clsPrefix:r,tmNode:e}):u(we,{clsPrefix:r,key:e.key,tmNode:e})))}),ce(e.action,e=>e&&[u(`div`,{class:`${r}-base-select-menu__action`,"data-action":!0,key:`action`},e),u(fe,{onFocus:this.onTabOut,key:`focus-detector`})]))}}),De=O([B(`base-selection`,`
 --n-padding-single: var(--n-padding-single-top) var(--n-padding-single-right) var(--n-padding-single-bottom) var(--n-padding-single-left);
 --n-padding-multiple: var(--n-padding-multiple-top) var(--n-padding-multiple-right) var(--n-padding-multiple-bottom) var(--n-padding-multiple-left);
 position: relative;
 z-index: auto;
 box-shadow: none;
 width: 100%;
 max-width: 100%;
 display: inline-block;
 vertical-align: bottom;
 border-radius: var(--n-border-radius);
 min-height: var(--n-height);
 line-height: 1.5;
 font-size: var(--n-font-size);
 `,[B(`base-loading`,`
 color: var(--n-loading-color);
 `),B(`base-selection-tags`,`min-height: var(--n-height);`),X(`border, state-border`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border: var(--n-border);
 border-radius: inherit;
 transition:
 box-shadow .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),X(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),B(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[X(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),B(`base-selection-overlay`,`
 display: flex;
 align-items: center;
 white-space: nowrap;
 pointer-events: none;
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 left: 0;
 padding: var(--n-padding-single);
 transition: color .3s var(--n-bezier);
 `,[X(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),B(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[X(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),B(`base-selection-tags`,`
 cursor: pointer;
 outline: none;
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 display: flex;
 padding: var(--n-padding-multiple);
 flex-wrap: wrap;
 align-items: center;
 width: 100%;
 vertical-align: bottom;
 background-color: var(--n-color);
 border-radius: inherit;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),B(`base-selection-label`,`
 height: var(--n-height);
 display: inline-flex;
 width: 100%;
 vertical-align: bottom;
 cursor: pointer;
 outline: none;
 z-index: auto;
 box-sizing: border-box;
 position: relative;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 border-radius: inherit;
 background-color: var(--n-color);
 align-items: center;
 `,[B(`base-selection-input`,`
 font-size: inherit;
 line-height: inherit;
 outline: none;
 cursor: pointer;
 box-sizing: border-box;
 border:none;
 width: 100%;
 padding: var(--n-padding-single);
 background-color: #0000;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 caret-color: var(--n-caret-color);
 `,[X(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),X(`render-label`,`
 color: var(--n-text-color);
 `)]),M(`disabled`,[O(`&:hover`,[X(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),A(`focus`,[X(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),A(`active`,[X(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),B(`base-selection-label`,`background-color: var(--n-color-active);`),B(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),A(`disabled`,`cursor: not-allowed;`,[X(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),B(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[B(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),X(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),B(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),B(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),B(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[X(`input`,`
 font-size: inherit;
 font-family: inherit;
 min-width: 1px;
 padding: 0;
 background-color: #0000;
 outline: none;
 border: none;
 max-width: 100%;
 overflow: hidden;
 width: 1em;
 line-height: inherit;
 cursor: pointer;
 color: var(--n-text-color);
 caret-color: var(--n-caret-color);
 `),X(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>A(`${e}-status`,[X(`state-border`,`border: var(--n-border-${e});`),M(`disabled`,[O(`&:hover`,[X(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),A(`active`,[X(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),B(`base-selection-label`,`background-color: var(--n-color-active-${e});`),B(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),A(`focus`,[X(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),B(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),B(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[O(`&:last-child`,`padding-right: 0;`),B(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[X(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),Oe=o({name:`InternalSelection`,props:Object.assign(Object.assign({},L.props),{clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function}),setup(t){let{mergedClsPrefixRef:n,mergedRtlRef:r}=l(t),i=_(`InternalSelection`,r,n),a=F(null),o=F(null),s=F(null),c=F(null),u=F(null),d=F(null),p=F(null),m=F(null),h=F(null),g=F(null),v=F(!1),y=F(!1),b=F(!1),x=L(`InternalSelection`,`-internal-selection`,De,me,t,w(t,`clsPrefix`)),C=f(()=>t.clearable&&!t.disabled&&(b.value||t.active)),ee=f(()=>t.selectedOption?t.renderTag?t.renderTag({option:t.selectedOption,handleClose:()=>{}}):t.renderLabel?t.renderLabel(t.selectedOption,!0):U(t.selectedOption[t.labelField],t.selectedOption,!0):t.placeholder),te=f(()=>{let e=t.selectedOption;if(e)return e[t.labelField]}),T=f(()=>t.multiple?!!(Array.isArray(t.selectedOptions)&&t.selectedOptions.length):t.selectedOption!==null);function E(){var e;let{value:n}=a;if(n){let{value:r}=o;r&&(r.style.width=`${n.offsetWidth}px`,t.maxTagCount!==`responsive`&&((e=h.value)==null||e.sync({showAllItemsBeforeCalculate:!1})))}}function O(){let{value:e}=g;e&&(e.style.display=`none`)}function k(){let{value:e}=g;e&&(e.style.display=`inline-block`)}Y(w(t,`active`),e=>{e||O()}),Y(w(t,`pattern`),()=>{t.multiple&&S(E)});function A(e){let{onFocus:n}=t;n&&n(e)}function ne(e){let{onBlur:n}=t;n&&n(e)}function j(e){let{onDeleteOption:n}=t;n&&n(e)}function M(e){let{onClear:n}=t;n&&n(e)}function N(e){let{onPatternInput:n}=t;n&&n(e)}function P(e){(!e.relatedTarget||!s.value?.contains(e.relatedTarget))&&A(e)}function I(e){s.value?.contains(e.relatedTarget)||ne(e)}function R(e){M(e)}function z(){b.value=!0}function B(){b.value=!1}function re(e){!t.active||!t.filterable||e.target!==o.value&&e.preventDefault()}function ie(e){j(e)}let V=F(!1);function H(e){if(e.key===`Backspace`&&!V.value&&!t.pattern.length){let{selectedOptions:e}=t;e?.length&&ie(e[e.length-1])}}let W=null;function G(e){let{value:n}=a;n&&(n.textContent=e.target.value,E()),t.ignoreComposition&&V.value?W=e:N(e)}function q(){V.value=!0}function J(){V.value=!1,t.ignoreComposition&&N(W),W=null}function se(e){var n;y.value=!0,(n=t.onPatternFocus)==null||n.call(t,e)}function ce(e){var n;y.value=!1,(n=t.onPatternBlur)==null||n.call(t,e)}function le(){var e,n;if(t.filterable)y.value=!1,(e=d.value)==null||e.blur(),(n=o.value)==null||n.blur();else if(t.multiple){let{value:e}=c;e?.blur()}else{let{value:e}=u;e?.blur()}}function X(){var e,n,r;t.filterable?(y.value=!1,(e=d.value)==null||e.focus()):t.multiple?(n=c.value)==null||n.focus():(r=u.value)==null||r.focus()}function ue(){let{value:e}=o;e&&(k(),e.focus())}function Z(){let{value:e}=o;e&&e.blur()}function Q(e){let{value:t}=p;t&&t.setTextContent(`+${e}`)}function de(){let{value:e}=m;return e}function fe(){return o.value}let pe=null;function he(){pe!==null&&window.clearTimeout(pe)}function ge(){t.active||(he(),pe=window.setTimeout(()=>{T.value&&(v.value=!0)},100))}function $(){he()}function _e(e){e||(he(),v.value=!1)}Y(T,e=>{e||(v.value=!1)}),e(()=>{D(()=>{let e=d.value;e&&(t.disabled?e.removeAttribute(`tabindex`):e.tabIndex=y.value?-1:0)})}),ye(s,t.onResize);let{inlineThemeDisabled:ve}=t,be=f(()=>{let{size:e}=t,{common:{cubicBezierEaseInOut:n},self:{fontWeight:r,borderRadius:i,color:a,placeholderColor:o,textColor:s,paddingSingle:c,paddingMultiple:l,caretColor:u,colorDisabled:d,textColorDisabled:f,placeholderColorDisabled:p,colorActive:m,boxShadowFocus:h,boxShadowActive:g,boxShadowHover:_,border:v,borderFocus:y,borderHover:b,borderActive:S,arrowColor:C,arrowColorDisabled:ee,loadingColor:w,colorActiveWarning:te,boxShadowFocusWarning:T,boxShadowActiveWarning:E,boxShadowHoverWarning:D,borderWarning:O,borderFocusWarning:k,borderHoverWarning:A,borderActiveWarning:ne,colorActiveError:j,boxShadowFocusError:M,boxShadowActiveError:N,boxShadowHoverError:P,borderError:F,borderFocusError:I,borderHoverError:L,borderActiveError:R,clearColor:z,clearColorHover:B,clearColorPressed:re,clearSize:ie,arrowSize:V,[K(`height`,e)]:H,[K(`fontSize`,e)]:U}}=x.value,W=oe(c),G=oe(l);return{"--n-bezier":n,"--n-border":v,"--n-border-active":S,"--n-border-focus":y,"--n-border-hover":b,"--n-border-radius":i,"--n-box-shadow-active":g,"--n-box-shadow-focus":h,"--n-box-shadow-hover":_,"--n-caret-color":u,"--n-color":a,"--n-color-active":m,"--n-color-disabled":d,"--n-font-size":U,"--n-height":H,"--n-padding-single-top":W.top,"--n-padding-multiple-top":G.top,"--n-padding-single-right":W.right,"--n-padding-multiple-right":G.right,"--n-padding-single-left":W.left,"--n-padding-multiple-left":G.left,"--n-padding-single-bottom":W.bottom,"--n-padding-multiple-bottom":G.bottom,"--n-placeholder-color":o,"--n-placeholder-color-disabled":p,"--n-text-color":s,"--n-text-color-disabled":f,"--n-arrow-color":C,"--n-arrow-color-disabled":ee,"--n-loading-color":w,"--n-color-active-warning":te,"--n-box-shadow-focus-warning":T,"--n-box-shadow-active-warning":E,"--n-box-shadow-hover-warning":D,"--n-border-warning":O,"--n-border-focus-warning":k,"--n-border-hover-warning":A,"--n-border-active-warning":ne,"--n-color-active-error":j,"--n-box-shadow-focus-error":M,"--n-box-shadow-active-error":N,"--n-box-shadow-hover-error":P,"--n-border-error":F,"--n-border-focus-error":I,"--n-border-hover-error":L,"--n-border-active-error":R,"--n-clear-size":ie,"--n-clear-color":z,"--n-clear-color-hover":B,"--n-clear-color-pressed":re,"--n-arrow-size":V,"--n-font-weight":r}}),xe=ve?ae(`internal-selection`,f(()=>t.size[0]),be,t):void 0;return{mergedTheme:x,mergedClearable:C,mergedClsPrefix:n,rtlEnabled:i,patternInputFocused:y,filterablePlaceholder:ee,label:te,selected:T,showTagsPanel:v,isComposing:V,counterRef:p,counterWrapperRef:m,patternInputMirrorRef:a,patternInputRef:o,selfRef:s,multipleElRef:c,singleElRef:u,patternInputWrapperRef:d,overflowRef:h,inputTagElRef:g,handleMouseDown:re,handleFocusin:P,handleClear:R,handleMouseEnter:z,handleMouseLeave:B,handleDeleteOption:ie,handlePatternKeyDown:H,handlePatternInputInput:G,handlePatternInputBlur:ce,handlePatternInputFocus:se,handleMouseEnterCounter:ge,handleMouseLeaveCounter:$,handleFocusout:I,handleCompositionEnd:J,handleCompositionStart:q,onPopoverUpdateShow:_e,focus:X,focusInput:ue,blur:le,blurInput:Z,updateCounter:Q,getCounter:de,getTail:fe,renderLabel:t.renderLabel,cssVars:ve?void 0:be,themeClass:xe?.themeClass,onRender:xe?.onRender}},render(){let{status:e,multiple:t,size:n,disabled:r,filterable:i,maxTagCount:o,bordered:s,clsPrefix:c,ellipsisTagPopoverProps:l,onRender:d,renderTag:f,renderLabel:p}=this;d?.();let m=o===`responsive`,h=typeof o==`number`,g=m||h,_=u(ee,null,{default:()=>u(j,{clsPrefix:c,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>{var e;return(e=this.$slots).arrow?.call(e)}})}),v;if(t){let{labelField:e}=this,t=t=>u(`div`,{class:`${c}-base-selection-tag-wrapper`,key:t.value},f?f({option:t,handleClose:()=>{this.handleDeleteOption(t)}}):u(ge,{size:n,closable:!t.disabled,disabled:r,onClose:()=>{this.handleDeleteOption(t)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>p?p(t,!0):U(t[e],t,!0)})),s=()=>(h?this.selectedOptions.slice(0,o):this.selectedOptions).map(t),d=i?u(`div`,{class:`${c}-base-selection-input-tag`,ref:`inputTagElRef`,key:`__input-tag__`},u(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:r,value:this.pattern,autofocus:this.autofocus,class:`${c}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),u(`span`,{ref:`patternInputMirrorRef`,class:`${c}-base-selection-input-tag__mirror`},this.pattern)):null,y=m?()=>u(`div`,{class:`${c}-base-selection-tag-wrapper`,ref:`counterWrapperRef`},u(ge,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:r})):void 0,b;if(h){let e=this.selectedOptions.length-o;e>0&&(b=u(`div`,{class:`${c}-base-selection-tag-wrapper`,key:`__counter__`},u(ge,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:r},{default:()=>`+${e}`})))}let x=m?i?u(ve,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:y,tail:()=>d}):u(ve,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:y}):h&&b?s().concat(b):s(),S=g?()=>u(`div`,{class:`${c}-base-selection-popover`},m?s():this.selectedOptions.map(t)):void 0,C=g?Object.assign({show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover},l):null,ee=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?u(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`},u(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):null,w=i?u(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-tags`},x,m?null:d,_):u(`div`,{ref:`multipleElRef`,class:`${c}-base-selection-tags`,tabindex:r?void 0:0},x,_);v=u(a,null,g?u(se,Object.assign({},C,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>w,default:S}):w,ee)}else if(i){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,n=this.active?!1:this.selected;v=u(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-label`,title:this.patternInputFocused?void 0:N(this.label)},u(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,class:`${c}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:r,disabled:r,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),n?u(`div`,{class:`${c}-base-selection-label__render-label ${c}-base-selection-overlay`,key:`input`},u(`div`,{class:`${c}-base-selection-overlay__wrapper`},f?f({option:this.selectedOption,handleClose:()=>{}}):p?p(this.selectedOption,!0):U(this.label,this.selectedOption,!0))):null,t?u(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},u(`div`,{class:`${c}-base-selection-overlay__wrapper`},this.filterablePlaceholder)):null,_)}else v=u(`div`,{ref:`singleElRef`,class:`${c}-base-selection-label`,tabindex:this.disabled?void 0:0},this.label===void 0?u(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},u(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):u(`div`,{class:`${c}-base-selection-input`,title:N(this.label),key:`input`},u(`div`,{class:`${c}-base-selection-input__content`},f?f({option:this.selectedOption,handleClose:()=>{}}):p?p(this.selectedOption,!0):U(this.label,this.selectedOption,!0))),_);return u(`div`,{ref:`selfRef`,class:[`${c}-base-selection`,this.rtlEnabled&&`${c}-base-selection--rtl`,this.themeClass,e&&`${c}-base-selection--${e}-status`,{[`${c}-base-selection--active`]:this.active,[`${c}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${c}-base-selection--disabled`]:this.disabled,[`${c}-base-selection--multiple`]:this.multiple,[`${c}-base-selection--focus`]:this.focused}],style:this.cssVars,onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},v,s?u(`div`,{class:`${c}-base-selection__border`}):null,s?u(`div`,{class:`${c}-base-selection__state-border`}):null)}});function ke(e){return e.type===`group`}function Ae(e){return e.type===`ignored`}function je(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function Me(e,t){return{getIsGroup:ke,getIgnored:Ae,getKey(t){return ke(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function Ne(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(ke(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if(Ae(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function Pe(e,t,n){let r=new Map;return e.forEach(e=>{ke(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}var Fe=O([B(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),B(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[T({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),Ie=o({name:`Select`,props:Object.assign(Object.assign({},L.props),{to:h.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array}),slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:n,namespaceRef:r,inlineThemeDisabled:i,mergedComponentPropsRef:a}=l(e),o=L(`Select`,`-select`,Fe,he,e,t),u=F(e.defaultValue),p=b(w(e,`value`),u),m=F(!1),g=F(``),_=x(e,[`items`,`options`]),v=F([]),y=F([]),S=f(()=>y.value.concat(v.value).concat(_.value)),C=f(()=>{let{filter:t}=e;if(t)return t;let{labelField:n,valueField:r}=e;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return je(e,i);let a=t[r];return typeof a==`string`?je(e,a):typeof a==`number`?je(e,String(a)):!1}}),ee=f(()=>{if(e.remote)return _.value;{let{value:t}=S,{value:n}=g;return!n.length||!e.filterable?t:Ne(t,C.value,n,e.childrenField)}}),T=f(()=>{let{valueField:t,childrenField:n}=e,r=Me(t,n);return re(ee.value,r)}),E=f(()=>Pe(S.value,e.valueField,e.childrenField)),D=F(!1),O=b(w(e,`show`),D),A=F(null),ne=F(null),j=F(null),{localeRef:M}=c(`Select`),N=f(()=>e.placeholder??M.value.placeholder),P=[],I=F(new Map),R=f(()=>{let{fallbackOption:t}=e;if(t===void 0){let{labelField:t,valueField:n}=e;return e=>({[t]:String(e),[n]:e})}return t===!1?!1:e=>Object.assign(t(e),{value:e})});function z(t){let n=e.remote,{value:r}=I,{value:i}=E,{value:a}=R,o=[];return t.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let B=f(()=>{if(e.multiple){let{value:e}=p;return Array.isArray(e)?z(e):[]}return null}),V=f(()=>{let{value:t}=p;return!e.multiple&&!Array.isArray(t)?t===null?null:z([t])[0]||null:null}),H=d(e,{mergedSize:t=>{let{size:n}=e;if(n)return n;let{mergedSize:r}=t||{};return r?.value?r.value:a?.value?.Select?.size||`medium`}}),{mergedSizeRef:U,mergedDisabledRef:G,mergedStatusRef:K}=H;function q(t,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=e,{nTriggerFormChange:o,nTriggerFormInput:s}=H;r&&k(r,t,n),a&&k(a,t,n),i&&k(i,t,n),u.value=t,o(),s()}function J(t){let{onBlur:n}=e,{nTriggerFormBlur:r}=H;n&&k(n,t),r()}function oe(){let{onClear:t}=e;t&&k(t)}function se(t){let{onFocus:n,showOnFocus:r}=e,{nTriggerFormFocus:i}=H;n&&k(n,t),i(),r&&Z()}function ce(t){let{onSearch:n}=e;n&&k(n,t)}function le(t){let{onScroll:n}=e;n&&k(n,t)}function X(){var t;let{remote:n,multiple:r}=e;if(n){let{value:n}=I;if(r){let{valueField:r}=e;(t=B.value)==null||t.forEach(e=>{n.set(e[r],e)})}else{let t=V.value;t&&n.set(t[e.valueField],t)}}}function ue(t){let{onUpdateShow:n,"onUpdate:show":r}=e;n&&k(n,t),r&&k(r,t),D.value=t}function Z(){G.value||(ue(!0),D.value=!0,e.filterable&&Le())}function Q(){ue(!1)}function de(){g.value=``,y.value=P}let fe=F(!1);function pe(){e.filterable&&(fe.value=!0)}function me(){e.filterable&&(fe.value=!1,O.value||de())}function ge(){G.value||(O.value?e.filterable?Le():Q():Z())}function $(e){(j.value?.selfRef)?.contains(e.relatedTarget)||(m.value=!1,J(e),Q())}function _e(e){se(e),m.value=!0}function ve(){m.value=!0}function ye(e){A.value?.$el.contains(e.relatedTarget)||(m.value=!1,J(e),Q())}function be(){var e;(e=A.value)==null||e.focus(),Q()}function xe(e){O.value&&(A.value?.$el.contains(te(e))||Q())}function Se(t){if(!Array.isArray(t))return[];if(R.value)return Array.from(t);{let{remote:n}=e,{value:r}=E;if(n){let{value:e}=I;return t.filter(t=>r.has(t)||e.has(t))}else return t.filter(e=>r.has(e))}}function Ce(e){we(e.rawNode)}function we(t){if(G.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=e;if(n&&!r){let{value:e}=y,t=e[0]||null;if(t){let e=v.value;e.length?e.push(t):v.value=[t],y.value=P}}if(r&&I.value.set(t[a],t),e.multiple){let e=Se(p.value),o=e.findIndex(e=>e===t[a]);if(~o){if(e.splice(o,1),n&&!r){let e=Te(t[a]);~e&&(v.value.splice(e,1),i&&(g.value=``))}}else e.push(t[a]),i&&(g.value=``);q(e,z(e))}else{if(n&&!r){let e=Te(t[a]);~e?v.value=[v.value[e]]:v.value=P}Ie(),Q(),q(t[a],t)}}function Te(t){return v.value.findIndex(n=>n[e.valueField]===t)}function Ee(t){O.value||Z();let{value:n}=t.target;g.value=n;let{tag:r,remote:i}=e;if(ce(n),r&&!i){if(!n){y.value=P;return}let{onCreate:t}=e,r=t?t(n):{[e.labelField]:n,[e.valueField]:n},{valueField:i,labelField:a}=e;_.value.some(e=>e[i]===r[i]||e[a]===r[a])||v.value.some(e=>e[i]===r[i]||e[a]===r[a])?y.value=P:y.value=[r]}}function De(t){t.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=e;!n&&e.filterable&&Q(),r&&!i&&a&&(v.value=P),oe(),n?q([],[]):q(null,null)}function Oe(e){!W(e,`action`)&&!W(e,`empty`)&&!W(e,`header`)&&e.preventDefault()}function ke(e){le(e)}function Ae(t){var n,r,i;if(!e.keyboard){t.preventDefault();return}switch(t.key){case` `:if(e.filterable)break;t.preventDefault();case`Enter`:if(!A.value?.isComposing){if(O.value){let t=j.value?.getPendingTmNode();t?Ce(t):e.filterable||(Q(),Ie())}else if(Z(),e.tag&&fe.value){let t=y.value[0];if(t){let n=t[e.valueField],{value:r}=p;e.multiple&&Array.isArray(r)&&r.includes(n)||we(t)}}}t.preventDefault();break;case`ArrowUp`:if(t.preventDefault(),e.loading)return;O.value&&((n=j.value)==null||n.prev());break;case`ArrowDown`:if(t.preventDefault(),e.loading)return;O.value?(r=j.value)==null||r.next():Z();break;case`Escape`:O.value&&(ie(t),Q()),(i=A.value)==null||i.focus();break}}function Ie(){var e;(e=A.value)==null||e.focus()}function Le(){var e;(e=A.value)==null||e.focusInput()}function Re(){var e;O.value&&((e=ne.value)==null||e.syncPosition())}X(),Y(w(e,`options`),X);let ze={focus:()=>{var e;(e=A.value)==null||e.focus()},focusInput:()=>{var e;(e=A.value)==null||e.focusInput()},blur:()=>{var e;(e=A.value)==null||e.blur()},blurInput:()=>{var e;(e=A.value)==null||e.blurInput()}},Be=f(()=>{let{self:{menuBoxShadow:e}}=o.value;return{"--n-menu-box-shadow":e}}),Ve=i?ae(`select`,void 0,Be,e):void 0;return Object.assign(Object.assign({},ze),{mergedStatus:K,mergedClsPrefix:t,mergedBordered:n,namespace:r,treeMate:T,isMounted:s(),triggerRef:A,menuRef:j,pattern:g,uncontrolledShow:D,mergedShow:O,adjustedTo:h(e),uncontrolledValue:u,mergedValue:p,followerRef:ne,localizedPlaceholder:N,selectedOption:V,selectedOptions:B,mergedSize:U,mergedDisabled:G,focused:m,activeWithoutMenuOpen:fe,inlineThemeDisabled:i,onTriggerInputFocus:pe,onTriggerInputBlur:me,handleTriggerOrMenuResize:Re,handleMenuFocus:ve,handleMenuBlur:ye,handleMenuTabOut:be,handleTriggerClick:ge,handleToggle:Ce,handleDeleteOption:we,handlePatternInput:Ee,handleClear:De,handleTriggerBlur:$,handleTriggerFocus:_e,handleKeydown:Ae,handleMenuAfterLeave:de,handleMenuClickOutside:xe,handleMenuScroll:ke,handleMenuKeydown:Ae,handleMenuMousedown:Oe,mergedTheme:o,cssVars:i?void 0:Be,themeClass:Ve?.themeClass,onRender:Ve?.onRender})},render(){return u(`div`,{class:`${this.mergedClsPrefix}-select`},u(g,null,{default:()=>[u(r,null,{default:()=>u(Oe,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{arrow:()=>{var e;return[(e=this.$slots).arrow?.call(e)]}})}),u(E,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===h.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{default:()=>u(R,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{default:()=>{var e;return this.mergedShow||this.displayDirective===`show`?((e=this.onRender)==null||e.call(this),le(u(Ee,Object.assign({},this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{empty:()=>{var e;return[(e=this.$slots).empty?.call(e)]},header:()=>{var e;return[(e=this.$slots).header?.call(e)]},action:()=>{var e;return[(e=this.$slots).action?.call(e)]}}),this.displayDirective===`show`?[[m,this.mergedShow],[i,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[i,this.handleMenuClickOutside,void 0,{capture:!0}]])):null}})})]}))}});export{be as i,Me as n,Ee as r,Ie as t};