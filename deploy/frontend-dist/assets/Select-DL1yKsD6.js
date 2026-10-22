import{$i as e,Ar as t,Bn as n,Cn as r,Di as i,En as a,Fr as o,Gi as s,Gn as c,Ji as l,Jn as u,Jr as d,Kt as f,Li as p,Nr as m,On as h,Pr as g,Qr as _,Qt as v,Ri as y,Rn as b,Ui as x,Un as S,Vi as C,Xi as w,Yr as T,Zr as E,Zt as D,ai as O,br as k,cn as A,dr as j,ea as ee,ei as M,er as N,fr as te,gr as ne,hr as P,ii as F,ir as I,ki as L,kr as R,la as z,mr as B,na as V,ni as H,nn as re,or as ie,qn as U,qt as W,ri as G,rn as K,sa as q,ti as J,tn as ae,vr as Y,wi as oe,wr as se,xi as ce,xn as le,yr as ue,zi as X,zn as de,zr as Z}from"./router-BEJlwQm6.js";import{o as fe,t as pe}from"./focus-detector-CHmIQMZQ.js";import{a as me,r as he,t as ge}from"./light-D4c-kMKx.js";import{n as _e}from"./data-freshness-DcN7nqj6.js";var Q=`v-hidden`,ve=B(`[v-hidden]`,{display:`none!important`}),ye=p({name:`Overflow`,props:{getCounter:Function,getTail:Function,updateCounter:Function,onUpdateCount:Function,onUpdateOverflow:Function},setup(e,{slots:t}){let n=q(null),r=q(null);function i(i){let{value:a}=n,{getCounter:o,getTail:s}=e,c;if(c=o===void 0?r.value:o(),!a||!c)return;c.hasAttribute(Q)&&c.removeAttribute(Q);let{children:l}=a;if(i.showAllItemsBeforeCalculate)for(let e of l)e.hasAttribute(Q)&&e.removeAttribute(Q);let u=a.offsetWidth,d=[],f=t.tail?s?.():null,p=f?f.offsetWidth:0,m=!1,h=a.children.length-+!!t.tail;for(let t=0;t<h-1;++t){if(t<0)continue;let n=l[t];if(m){n.hasAttribute(Q)||n.setAttribute(Q,``);continue}else n.hasAttribute(Q)&&n.removeAttribute(Q);let r=n.offsetWidth;if(p+=r,d[t]=r,p>u){let{updateCounter:n}=e;for(let r=t;r>=0;--r){let i=h-1-r;n===void 0?c.textContent=`${i}`:n(i);let a=c.offsetWidth;if(p-=d[r],p+a<=u||r===0){m=!0,t=r-1,f&&(t===-1?(f.style.maxWidth=`${u-a}px`,f.style.boxSizing=`border-box`):f.style.maxWidth=``);let{onUpdateCount:n}=e;n&&n(i);break}}}}let{onUpdateOverflow:g}=e;m?g!==void 0&&g(!0):(g!==void 0&&g(!1),c.setAttribute(Q,``))}let a=ne();return ve.mount({id:`vueuc/overflow`,head:!0,anchorMetaName:P,ssr:a}),s(()=>i({showAllItemsBeforeCalculate:!1})),{selfRef:n,counterRef:r,sync:i}},render(){let{$slots:e}=this;return C(()=>this.sync({showAllItemsBeforeCalculate:!1})),y(`div`,{class:`v-overflow`,ref:`selfRef`},[w(e,`default`),e.counter?e.counter():y(`span`,{style:{display:`inline-block`},ref:`counterRef`}),e.tail?e.tail():null])}});function be(t,n){n&&(s(()=>{let{value:e}=t;e&&j.registerHandler(e,n)}),e(t,(e,t)=>{t&&j.unregisterHandler(t)},{deep:!1}),x(()=>{let{value:e}=t;e&&j.unregisterHandler(e)}))}function $(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var xe=p({name:`Checkmark`,render(){return y(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},y(`g`,{fill:`none`},y(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})))}}),Se=p({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:n,labelFieldRef:r,nodePropsRef:i}=X(t);return{labelField:r,nodeProps:i,renderLabel:e,renderOption:n}},render(){let{clsPrefix:e,renderLabel:t,renderOption:n,nodeProps:r,tmNode:{rawNode:i}}=this,a=r?.(i),o=t?t(i,!1):u(i[this.labelField],i,!1),s=y(`div`,Object.assign({},a,{class:[`${e}-base-select-group-header`,a?.class]}),o);return i.render?i.render({node:s,option:i}):n?n({node:s,option:i,selected:!1}):s}});function Ce(e,t){return y(ce,{name:`fade-in-scale-up-transition`},{default:()=>e?y(le,{clsPrefix:t,class:`${t}-base-select-option__check`},{default:()=>y(xe)}):null})}var we=p({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:n,pendingTmNodeRef:r,multipleRef:i,valueSetRef:a,renderLabelRef:o,renderOptionRef:s,labelFieldRef:c,valueFieldRef:l,showCheckmarkRef:u,nodePropsRef:d,handleOptionClick:f,handleOptionMouseEnter:p}=X(t),m=Z(()=>{let{value:t}=r;return t?e.tmNode.key===t.key:!1});function h(t){let{tmNode:n}=e;n.disabled||f(t,n)}function g(t){let{tmNode:n}=e;n.disabled||p(t,n)}function _(t){let{tmNode:n}=e,{value:r}=m;n.disabled||r||p(t,n)}return{multiple:i,isGrouped:Z(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:u,nodeProps:d,isPending:m,isSelected:Z(()=>{let{value:t}=n,{value:r}=i;if(t===null)return!1;let o=e.tmNode.rawNode[l.value];if(r){let{value:e}=a;return e.has(o)}else return t===o}),labelField:c,renderLabel:o,renderOption:s,handleMouseMove:_,handleMouseEnter:g,handleClick:h}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:n,isPending:r,isGrouped:i,showCheckmark:a,nodeProps:o,renderOption:s,renderLabel:c,handleClick:l,handleMouseEnter:d,handleMouseMove:f}=this,p=Ce(n,e),m=c?[c(t,n),a&&p]:[u(t[this.labelField],t,n),a&&p],h=o?.(t),g=y(`div`,Object.assign({},h,{class:[`${e}-base-select-option`,t.class,h?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:n,[`${e}-base-select-option--grouped`]:i,[`${e}-base-select-option--pending`]:r,[`${e}-base-select-option--show-checkmark`]:a}],style:[h?.style||``,t.style||``],onClick:$([l,h?.onClick]),onMouseenter:$([d,h?.onMouseenter]),onMousemove:$([f,h?.onMousemove])}),y(`div`,{class:`${e}-base-select-option__content`},m));return t.render?t.render({node:g,option:t,selected:n}):s?s({node:g,option:t,selected:n}):g}}),Te=J(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[J(`scrollbar`,`
 max-height: var(--n-height);
 `),J(`virtual-list`,`
 max-height: var(--n-height);
 `),J(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[H(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),J(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),J(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),H(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),H(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),H(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),H(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),J(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),J(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[G(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),M(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),M(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),G(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),G(`pending`,[M(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),G(`selected`,`
 color: var(--n-option-text-color-active);
 `,[M(`&::before`,`
 background-color: var(--n-option-color-active);
 `),G(`pending`,[M(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),G(`disabled`,`
 cursor: not-allowed;
 `,[F(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),G(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),H(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[D({enterScale:`0.5`})])])]),Ee=p({name:`InternalSelectMenu`,props:Object.assign(Object.assign({},r.props),{clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function}),setup(i){let{mergedClsPrefixRef:o,mergedRtlRef:c,mergedComponentPropsRef:u}=n(i),f=a(`InternalSelectMenu`,c,o),p=r(`InternalSelectMenu`,`-internal-select-menu`,Te,me,i,z(i,`clsPrefix`)),m=q(null),h=q(null),g=q(null),v=L(()=>i.treeMate.getFlattenedNodes()),y=L(()=>re(v.value)),b=q(null);function S(){let{treeMate:e}=i,t=null,{value:n}=i;n===null?t=e.getFirstAvailableNode():(t=i.multiple?e.getNode((n||[])[(n||[]).length-1]):e.getNode(n),(!t||t.disabled)&&(t=e.getFirstAvailableNode())),W(t||null)}function w(){let{value:e}=b;e&&!i.treeMate.getNode(e.key)&&(b.value=null)}let E;e(()=>i.show,t=>{t?E=e(()=>i.treeMate,()=>{i.resetMenuOnOptionsChange?(i.autoPending?S():w(),C(G)):w()},{immediate:!0}):E?.()},{immediate:!0}),x(()=>{E?.()});let D=L(()=>d(p.value.self[O(`optionHeight`,i.size)])),k=L(()=>T(p.value.self[O(`padding`,i.size)])),A=L(()=>i.multiple&&Array.isArray(i.value)?new Set(i.value):new Set),j=L(()=>{let e=v.value;return e&&e.length===0}),ee=L(()=>u?.value?.Select?.renderEmpty);function M(e){let{onToggle:t}=i;t&&t(e)}function N(e){let{onScroll:t}=i;t&&t(e)}function te(e){var t;(t=g.value)==null||t.sync(),N(e)}function ne(){var e;(e=g.value)==null||e.sync()}function P(){let{value:e}=b;return e||null}function F(e,t){t.disabled||W(t,!1)}function I(e,t){t.disabled||M(t)}function B(e){var t;_(e,`action`)||(t=i.onKeyup)==null||t.call(i,e)}function V(e){var t;_(e,`action`)||(t=i.onKeydown)==null||t.call(i,e)}function H(e){var t;(t=i.onMousedown)==null||t.call(i,e),!i.focusable&&e.preventDefault()}function ie(){let{value:e}=b;e&&W(e.getNext({loop:!0}),!0)}function U(){let{value:e}=b;e&&W(e.getPrev({loop:!0}),!0)}function W(e,t=!1){b.value=e,t&&G()}function G(){var e,t;let n=b.value;if(!n)return;let r=y.value(n.key);r!==null&&(i.virtualScroll?(e=h.value)==null||e.scrollTo({index:r}):(t=g.value)==null||t.scrollTo({index:r,elSize:D.value}))}function K(e){var t;m.value?.contains(e.target)&&((t=i.onFocus)==null||t.call(i,e))}function J(e){var t;m.value?.contains(e.relatedTarget)||(t=i.onBlur)==null||t.call(i,e)}l(t,{handleOptionMouseEnter:F,handleOptionClick:I,valueSetRef:A,pendingTmNodeRef:b,nodePropsRef:z(i,`nodeProps`),showCheckmarkRef:z(i,`showCheckmark`),multipleRef:z(i,`multiple`),valueRef:z(i,`value`),renderLabelRef:z(i,`renderLabel`),renderOptionRef:z(i,`renderOption`),labelFieldRef:z(i,`labelField`),valueFieldRef:z(i,`valueField`)}),l(R,m),s(()=>{let{value:e}=g;e&&e.sync()});let ae=L(()=>{let{size:e}=i,{common:{cubicBezierEaseInOut:t},self:{height:n,borderRadius:r,color:a,groupHeaderTextColor:o,actionDividerColor:s,optionTextColorPressed:c,optionTextColor:l,optionTextColorDisabled:u,optionTextColorActive:d,optionOpacityDisabled:f,optionCheckColor:m,actionTextColor:h,optionColorPending:g,optionColorActive:_,loadingColor:v,loadingSize:y,optionColorActivePending:b,[O(`optionFontSize`,e)]:x,[O(`optionHeight`,e)]:S,[O(`optionPadding`,e)]:C}}=p.value;return{"--n-height":n,"--n-action-divider-color":s,"--n-action-text-color":h,"--n-bezier":t,"--n-border-radius":r,"--n-color":a,"--n-option-font-size":x,"--n-group-header-text-color":o,"--n-option-check-color":m,"--n-option-color-pending":g,"--n-option-color-active":_,"--n-option-color-active-pending":b,"--n-option-height":S,"--n-option-opacity-disabled":f,"--n-option-text-color":l,"--n-option-text-color-active":d,"--n-option-text-color-disabled":u,"--n-option-text-color-pressed":c,"--n-option-padding":C,"--n-option-padding-left":T(C,`left`),"--n-option-padding-right":T(C,`right`),"--n-loading-color":v,"--n-loading-size":y}}),{inlineThemeDisabled:Y}=i,oe=Y?de(`internal-select-menu`,L(()=>i.size[0]),ae,i):void 0,se={selfRef:m,next:ie,prev:U,getPendingTmNode:P};return be(m,i.onResize),Object.assign({mergedTheme:p,mergedClsPrefix:o,rtlEnabled:f,virtualListRef:h,scrollbarRef:g,itemSize:D,padding:k,flattenedNodes:v,empty:j,mergedRenderEmpty:ee,virtualListContainer(){let{value:e}=h;return e?.listElRef},virtualListContent(){let{value:e}=h;return e?.itemsElRef},doScroll:N,handleFocusin:K,handleFocusout:J,handleKeyUp:B,handleKeyDown:V,handleMouseDown:H,handleVirtualListResize:ne,handleVirtualListScroll:te,cssVars:Y?void 0:ae,themeClass:oe?.themeClass,onRender:oe?.onRender},se)},render(){let{$slots:e,virtualScroll:t,clsPrefix:n,mergedTheme:r,themeClass:i,onRender:a}=this;return a?.(),y(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:[`${n}-base-select-menu`,`${n}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${n}-base-select-menu--rtl`,i,this.multiple&&`${n}-base-select-menu--multiple`],style:this.cssVars,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},U(e.header,e=>e&&y(`div`,{class:`${n}-base-select-menu__header`,"data-header":!0,key:`header`},e)),this.loading?y(`div`,{class:`${n}-base-select-menu__loading`},y(A,{clsPrefix:n,strokeWidth:20})):this.empty?y(`div`,{class:`${n}-base-select-menu__empty`,"data-empty":!0},c(e.empty,()=>[this.mergedRenderEmpty?.call(this)||y(v,{theme:r.peers.Empty,themeOverrides:r.peerOverrides.Empty,size:this.size})])):y(K,Object.assign({ref:`scrollbarRef`,theme:r.peers.Scrollbar,themeOverrides:r.peerOverrides.Scrollbar,scrollable:this.scrollable,container:t?this.virtualListContainer:void 0,content:t?this.virtualListContent:void 0,onScroll:t?void 0:this.doScroll},this.scrollbarProps),{default:()=>t?y(fe,{ref:`virtualListRef`,class:`${n}-virtual-list`,items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?y(Se,{key:e.key,clsPrefix:n,tmNode:e}):e.ignored?null:y(we,{clsPrefix:n,key:e.key,tmNode:e})}):y(`div`,{class:`${n}-base-select-menu-option-wrapper`,style:{paddingTop:this.padding.top,paddingBottom:this.padding.bottom}},this.flattenedNodes.map(e=>e.isGroup?y(Se,{key:e.key,clsPrefix:n,tmNode:e}):y(we,{clsPrefix:n,key:e.key,tmNode:e})))}),U(e.action,e=>e&&[y(`div`,{class:`${n}-base-select-menu__action`,"data-action":!0,key:`action`},e),y(pe,{onFocus:this.onTabOut,key:`focus-detector`})]))}}),De=M([J(`base-selection`,`
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
 `,[J(`base-loading`,`
 color: var(--n-loading-color);
 `),J(`base-selection-tags`,`min-height: var(--n-height);`),H(`border, state-border`,`
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
 `),H(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),J(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[H(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),J(`base-selection-overlay`,`
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
 `,[H(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),J(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[H(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),J(`base-selection-tags`,`
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
 `),J(`base-selection-label`,`
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
 `,[J(`base-selection-input`,`
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
 `,[H(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),H(`render-label`,`
 color: var(--n-text-color);
 `)]),F(`disabled`,[M(`&:hover`,[H(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),G(`focus`,[H(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),G(`active`,[H(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),J(`base-selection-label`,`background-color: var(--n-color-active);`),J(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),G(`disabled`,`cursor: not-allowed;`,[H(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),J(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[J(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),H(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),J(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),J(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),J(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[H(`input`,`
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
 `),H(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>G(`${e}-status`,[H(`state-border`,`border: var(--n-border-${e});`),F(`disabled`,[M(`&:hover`,[H(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),G(`active`,[H(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),J(`base-selection-label`,`background-color: var(--n-color-active-${e});`),J(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),G(`focus`,[H(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),J(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),J(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[M(`&:last-child`,`padding-right: 0;`),J(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[H(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),Oe=p({name:`InternalSelection`,props:Object.assign(Object.assign({},r.props),{clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function}),setup(t){let{mergedClsPrefixRef:i,mergedRtlRef:o}=n(t),c=a(`InternalSelection`,o,i),l=q(null),d=q(null),f=q(null),p=q(null),m=q(null),h=q(null),g=q(null),_=q(null),v=q(null),y=q(null),b=q(!1),x=q(!1),S=q(!1),w=r(`InternalSelection`,`-internal-selection`,De,he,t,z(t,`clsPrefix`)),E=L(()=>t.clearable&&!t.disabled&&(S.value||t.active)),D=L(()=>t.selectedOption?t.renderTag?t.renderTag({option:t.selectedOption,handleClose:()=>{}}):t.renderLabel?t.renderLabel(t.selectedOption,!0):u(t.selectedOption[t.labelField],t.selectedOption,!0):t.placeholder),k=L(()=>{let e=t.selectedOption;if(e)return e[t.labelField]}),A=L(()=>t.multiple?!!(Array.isArray(t.selectedOptions)&&t.selectedOptions.length):t.selectedOption!==null);function j(){var e;let{value:n}=l;if(n){let{value:r}=d;r&&(r.style.width=`${n.offsetWidth}px`,t.maxTagCount!==`responsive`&&((e=v.value)==null||e.sync({showAllItemsBeforeCalculate:!1})))}}function M(){let{value:e}=y;e&&(e.style.display=`none`)}function N(){let{value:e}=y;e&&(e.style.display=`inline-block`)}e(z(t,`active`),e=>{e||M()}),e(z(t,`pattern`),()=>{t.multiple&&C(j)});function te(e){let{onFocus:n}=t;n&&n(e)}function ne(e){let{onBlur:n}=t;n&&n(e)}function P(e){let{onDeleteOption:n}=t;n&&n(e)}function F(e){let{onClear:n}=t;n&&n(e)}function I(e){let{onPatternInput:n}=t;n&&n(e)}function R(e){(!e.relatedTarget||!f.value?.contains(e.relatedTarget))&&te(e)}function B(e){f.value?.contains(e.relatedTarget)||ne(e)}function V(e){F(e)}function H(){S.value=!0}function re(){S.value=!1}function ie(e){!t.active||!t.filterable||e.target!==d.value&&e.preventDefault()}function U(e){P(e)}let W=q(!1);function G(e){if(e.key===`Backspace`&&!W.value&&!t.pattern.length){let{selectedOptions:e}=t;e?.length&&U(e[e.length-1])}}let K=null;function J(e){let{value:n}=l;n&&(n.textContent=e.target.value,j()),t.ignoreComposition&&W.value?K=e:I(e)}function ae(){W.value=!0}function Y(){W.value=!1,t.ignoreComposition&&I(K),K=null}function oe(e){var n;x.value=!0,(n=t.onPatternFocus)==null||n.call(t,e)}function se(e){var n;x.value=!1,(n=t.onPatternBlur)==null||n.call(t,e)}function ce(){var e,n;if(t.filterable)x.value=!1,(e=h.value)==null||e.blur(),(n=d.value)==null||n.blur();else if(t.multiple){let{value:e}=p;e?.blur()}else{let{value:e}=m;e?.blur()}}function le(){var e,n,r;t.filterable?(x.value=!1,(e=h.value)==null||e.focus()):t.multiple?(n=p.value)==null||n.focus():(r=m.value)==null||r.focus()}function ue(){let{value:e}=d;e&&(N(),e.focus())}function X(){let{value:e}=d;e&&e.blur()}function Z(e){let{value:t}=g;t&&t.setTextContent(`+${e}`)}function fe(){let{value:e}=_;return e}function pe(){return d.value}let me=null;function ge(){me!==null&&window.clearTimeout(me)}function _e(){t.active||(ge(),me=window.setTimeout(()=>{A.value&&(b.value=!0)},100))}function Q(){ge()}function ve(e){e||(ge(),b.value=!1)}e(A,e=>{e||(b.value=!1)}),s(()=>{ee(()=>{let e=h.value;e&&(t.disabled?e.removeAttribute(`tabindex`):e.tabIndex=x.value?-1:0)})}),be(f,t.onResize);let{inlineThemeDisabled:ye}=t,$=L(()=>{let{size:e}=t,{common:{cubicBezierEaseInOut:n},self:{fontWeight:r,borderRadius:i,color:a,placeholderColor:o,textColor:s,paddingSingle:c,paddingMultiple:l,caretColor:u,colorDisabled:d,textColorDisabled:f,placeholderColorDisabled:p,colorActive:m,boxShadowFocus:h,boxShadowActive:g,boxShadowHover:_,border:v,borderFocus:y,borderHover:b,borderActive:x,arrowColor:S,arrowColorDisabled:C,loadingColor:E,colorActiveWarning:D,boxShadowFocusWarning:k,boxShadowActiveWarning:A,boxShadowHoverWarning:j,borderWarning:ee,borderFocusWarning:M,borderHoverWarning:N,borderActiveWarning:te,colorActiveError:ne,boxShadowFocusError:P,boxShadowActiveError:F,boxShadowHoverError:I,borderError:L,borderFocusError:R,borderHoverError:z,borderActiveError:B,clearColor:V,clearColorHover:H,clearColorPressed:re,clearSize:ie,arrowSize:U,[O(`height`,e)]:W,[O(`fontSize`,e)]:G}}=w.value,K=T(c),q=T(l);return{"--n-bezier":n,"--n-border":v,"--n-border-active":x,"--n-border-focus":y,"--n-border-hover":b,"--n-border-radius":i,"--n-box-shadow-active":g,"--n-box-shadow-focus":h,"--n-box-shadow-hover":_,"--n-caret-color":u,"--n-color":a,"--n-color-active":m,"--n-color-disabled":d,"--n-font-size":G,"--n-height":W,"--n-padding-single-top":K.top,"--n-padding-multiple-top":q.top,"--n-padding-single-right":K.right,"--n-padding-multiple-right":q.right,"--n-padding-single-left":K.left,"--n-padding-multiple-left":q.left,"--n-padding-single-bottom":K.bottom,"--n-padding-multiple-bottom":q.bottom,"--n-placeholder-color":o,"--n-placeholder-color-disabled":p,"--n-text-color":s,"--n-text-color-disabled":f,"--n-arrow-color":S,"--n-arrow-color-disabled":C,"--n-loading-color":E,"--n-color-active-warning":D,"--n-box-shadow-focus-warning":k,"--n-box-shadow-active-warning":A,"--n-box-shadow-hover-warning":j,"--n-border-warning":ee,"--n-border-focus-warning":M,"--n-border-hover-warning":N,"--n-border-active-warning":te,"--n-color-active-error":ne,"--n-box-shadow-focus-error":P,"--n-box-shadow-active-error":F,"--n-box-shadow-hover-error":I,"--n-border-error":L,"--n-border-focus-error":R,"--n-border-hover-error":z,"--n-border-active-error":B,"--n-clear-size":ie,"--n-clear-color":V,"--n-clear-color-hover":H,"--n-clear-color-pressed":re,"--n-arrow-size":U,"--n-font-weight":r}}),xe=ye?de(`internal-selection`,L(()=>t.size[0]),$,t):void 0;return{mergedTheme:w,mergedClearable:E,mergedClsPrefix:i,rtlEnabled:c,patternInputFocused:x,filterablePlaceholder:D,label:k,selected:A,showTagsPanel:b,isComposing:W,counterRef:g,counterWrapperRef:_,patternInputMirrorRef:l,patternInputRef:d,selfRef:f,multipleElRef:p,singleElRef:m,patternInputWrapperRef:h,overflowRef:v,inputTagElRef:y,handleMouseDown:ie,handleFocusin:R,handleClear:V,handleMouseEnter:H,handleMouseLeave:re,handleDeleteOption:U,handlePatternKeyDown:G,handlePatternInputInput:J,handlePatternInputBlur:se,handlePatternInputFocus:oe,handleMouseEnterCounter:_e,handleMouseLeaveCounter:Q,handleFocusout:B,handleCompositionEnd:Y,handleCompositionStart:ae,onPopoverUpdateShow:ve,focus:le,focusInput:ue,blur:ce,blurInput:X,updateCounter:Z,getCounter:fe,getTail:pe,renderLabel:t.renderLabel,cssVars:ye?void 0:$,themeClass:xe?.themeClass,onRender:xe?.onRender}},render(){let{status:e,multiple:t,size:n,disabled:r,filterable:a,maxTagCount:o,bordered:s,clsPrefix:c,ellipsisTagPopoverProps:l,onRender:d,renderTag:p,renderLabel:m}=this;d?.();let h=o===`responsive`,g=typeof o==`number`,_=h||g,v=y(S,null,{default:()=>y(f,{clsPrefix:c,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>{var e;return(e=this.$slots).arrow?.call(e)}})}),b;if(t){let{labelField:e}=this,t=t=>y(`div`,{class:`${c}-base-selection-tag-wrapper`,key:t.value},p?p({option:t,handleClose:()=>{this.handleDeleteOption(t)}}):y(_e,{size:n,closable:!t.disabled,disabled:r,onClose:()=>{this.handleDeleteOption(t)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>m?m(t,!0):u(t[e],t,!0)})),s=()=>(g?this.selectedOptions.slice(0,o):this.selectedOptions).map(t),d=a?y(`div`,{class:`${c}-base-selection-input-tag`,ref:`inputTagElRef`,key:`__input-tag__`},y(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:r,value:this.pattern,autofocus:this.autofocus,class:`${c}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),y(`span`,{ref:`patternInputMirrorRef`,class:`${c}-base-selection-input-tag__mirror`},this.pattern)):null,f=h?()=>y(`div`,{class:`${c}-base-selection-tag-wrapper`,ref:`counterWrapperRef`},y(_e,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:r})):void 0,x;if(g){let e=this.selectedOptions.length-o;e>0&&(x=y(`div`,{class:`${c}-base-selection-tag-wrapper`,key:`__counter__`},y(_e,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:r},{default:()=>`+${e}`})))}let S=h?a?y(ye,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:f,tail:()=>d}):y(ye,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:f}):g&&x?s().concat(x):s(),C=_?()=>y(`div`,{class:`${c}-base-selection-popover`},h?s():this.selectedOptions.map(t)):void 0,w=_?Object.assign({show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover},l):null,T=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?y(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`},y(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):null,E=a?y(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-tags`},S,h?null:d,v):y(`div`,{ref:`multipleElRef`,class:`${c}-base-selection-tags`,tabindex:r?void 0:0},S,v);b=y(i,null,_?y(W,Object.assign({},w,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>E,default:C}):E,T)}else if(a){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,n=this.active?!1:this.selected;b=y(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-label`,title:this.patternInputFocused?void 0:I(this.label)},y(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,class:`${c}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:r,disabled:r,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),n?y(`div`,{class:`${c}-base-selection-label__render-label ${c}-base-selection-overlay`,key:`input`},y(`div`,{class:`${c}-base-selection-overlay__wrapper`},p?p({option:this.selectedOption,handleClose:()=>{}}):m?m(this.selectedOption,!0):u(this.label,this.selectedOption,!0))):null,t?y(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},y(`div`,{class:`${c}-base-selection-overlay__wrapper`},this.filterablePlaceholder)):null,v)}else b=y(`div`,{ref:`singleElRef`,class:`${c}-base-selection-label`,tabindex:this.disabled?void 0:0},this.label===void 0?y(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},y(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):y(`div`,{class:`${c}-base-selection-input`,title:I(this.label),key:`input`},y(`div`,{class:`${c}-base-selection-input__content`},p?p({option:this.selectedOption,handleClose:()=>{}}):m?m(this.selectedOption,!0):u(this.label,this.selectedOption,!0))),v);return y(`div`,{ref:`selfRef`,class:[`${c}-base-selection`,this.rtlEnabled&&`${c}-base-selection--rtl`,this.themeClass,e&&`${c}-base-selection--${e}-status`,{[`${c}-base-selection--active`]:this.active,[`${c}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${c}-base-selection--disabled`]:this.disabled,[`${c}-base-selection--multiple`]:this.multiple,[`${c}-base-selection--focus`]:this.focused}],style:this.cssVars,onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},b,s?y(`div`,{class:`${c}-base-selection__border`}):null,s?y(`div`,{class:`${c}-base-selection__state-border`}):null)}});function ke(e){return e.type===`group`}function Ae(e){return e.type===`ignored`}function je(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function Me(e,t){return{getIsGroup:ke,getIgnored:Ae,getKey(t){return ke(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function Ne(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(ke(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if(Ae(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function Pe(e,t,n){let r=new Map;return e.forEach(e=>{ke(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}var Fe=M([J(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),J(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[D({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),Ie=p({name:`Select`,props:Object.assign(Object.assign({},r.props),{to:se.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array}),slots:Object,setup(t){let{mergedClsPrefixRef:i,mergedBorderedRef:a,namespaceRef:s,inlineThemeDisabled:c,mergedComponentPropsRef:l}=n(t),u=r(`Select`,`-select`,Fe,ge,t,i),d=q(t.defaultValue),f=o(z(t,`value`),d),p=q(!1),v=q(``),y=m(t,[`items`,`options`]),x=q([]),S=q([]),C=L(()=>S.value.concat(x.value).concat(y.value)),w=L(()=>{let{filter:e}=t;if(e)return e;let{labelField:n,valueField:r}=t;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return je(e,i);let a=t[r];return typeof a==`string`?je(e,a):typeof a==`number`?je(e,String(a)):!1}}),T=L(()=>{if(t.remote)return y.value;{let{value:e}=C,{value:n}=v;return!n.length||!t.filterable?e:Ne(e,w.value,n,t.childrenField)}}),D=L(()=>{let{valueField:e,childrenField:n}=t,r=Me(e,n);return ae(T.value,r)}),O=L(()=>Pe(C.value,t.valueField,t.childrenField)),k=q(!1),A=o(z(t,`show`),k),j=q(null),ee=q(null),M=q(null),{localeRef:te}=h(`Select`),ne=L(()=>t.placeholder??te.value.placeholder),P=[],F=q(new Map),I=L(()=>{let{fallbackOption:e}=t;if(e===void 0){let{labelField:e,valueField:n}=t;return t=>({[e]:String(t),[n]:t})}return e===!1?!1:t=>Object.assign(e(t),{value:t})});function R(e){let n=t.remote,{value:r}=F,{value:i}=O,{value:a}=I,o=[];return e.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let B=L(()=>{if(t.multiple){let{value:e}=f;return Array.isArray(e)?R(e):[]}return null}),V=L(()=>{let{value:e}=f;return!t.multiple&&!Array.isArray(e)?e===null?null:R([e])[0]||null:null}),H=b(t,{mergedSize:e=>{let{size:n}=t;if(n)return n;let{mergedSize:r}=e||{};return r?.value?r.value:l?.value?.Select?.size||`medium`}}),{mergedSizeRef:re,mergedDisabledRef:U,mergedStatusRef:W}=H;function G(e,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=t,{nTriggerFormChange:o,nTriggerFormInput:s}=H;r&&N(r,e,n),a&&N(a,e,n),i&&N(i,e,n),d.value=e,o(),s()}function K(e){let{onBlur:n}=t,{nTriggerFormBlur:r}=H;n&&N(n,e),r()}function J(){let{onClear:e}=t;e&&N(e)}function Y(e){let{onFocus:n,showOnFocus:r}=t,{nTriggerFormFocus:i}=H;n&&N(n,e),i(),r&&X()}function oe(e){let{onSearch:n}=t;n&&N(n,e)}function ce(e){let{onScroll:n}=t;n&&N(n,e)}function le(){var e;let{remote:n,multiple:r}=t;if(n){let{value:n}=F;if(r){let{valueField:r}=t;(e=B.value)==null||e.forEach(e=>{n.set(e[r],e)})}else{let e=V.value;e&&n.set(e[t.valueField],e)}}}function ue(e){let{onUpdateShow:n,"onUpdate:show":r}=t;n&&N(n,e),r&&N(r,e),k.value=e}function X(){U.value||(ue(!0),k.value=!0,t.filterable&&Le())}function Z(){ue(!1)}function fe(){v.value=``,S.value=P}let pe=q(!1);function me(){t.filterable&&(pe.value=!0)}function he(){t.filterable&&(pe.value=!1,A.value||fe())}function _e(){U.value||(A.value?t.filterable?Le():Z():X())}function Q(e){(M.value?.selfRef)?.contains(e.relatedTarget)||(p.value=!1,K(e),Z())}function ve(e){Y(e),p.value=!0}function ye(){p.value=!0}function be(e){j.value?.$el.contains(e.relatedTarget)||(p.value=!1,K(e),Z())}function $(){var e;(e=j.value)==null||e.focus(),Z()}function xe(e){A.value&&(j.value?.$el.contains(E(e))||Z())}function Se(e){if(!Array.isArray(e))return[];if(I.value)return Array.from(e);{let{remote:n}=t,{value:r}=O;if(n){let{value:t}=F;return e.filter(e=>r.has(e)||t.has(e))}else return e.filter(e=>r.has(e))}}function Ce(e){we(e.rawNode)}function we(e){if(U.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=t;if(n&&!r){let{value:e}=S,t=e[0]||null;if(t){let e=x.value;e.length?e.push(t):x.value=[t],S.value=P}}if(r&&F.value.set(e[a],e),t.multiple){let t=Se(f.value),o=t.findIndex(t=>t===e[a]);if(~o){if(t.splice(o,1),n&&!r){let t=Te(e[a]);~t&&(x.value.splice(t,1),i&&(v.value=``))}}else t.push(e[a]),i&&(v.value=``);G(t,R(t))}else{if(n&&!r){let t=Te(e[a]);~t?x.value=[x.value[t]]:x.value=P}Ie(),Z(),G(e[a],e)}}function Te(e){return x.value.findIndex(n=>n[t.valueField]===e)}function Ee(e){A.value||X();let{value:n}=e.target;v.value=n;let{tag:r,remote:i}=t;if(oe(n),r&&!i){if(!n){S.value=P;return}let{onCreate:e}=t,r=e?e(n):{[t.labelField]:n,[t.valueField]:n},{valueField:i,labelField:a}=t;y.value.some(e=>e[i]===r[i]||e[a]===r[a])||x.value.some(e=>e[i]===r[i]||e[a]===r[a])?S.value=P:S.value=[r]}}function De(e){e.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=t;!n&&t.filterable&&Z(),r&&!i&&a&&(x.value=P),J(),n?G([],[]):G(null,null)}function Oe(e){!_(e,`action`)&&!_(e,`empty`)&&!_(e,`header`)&&e.preventDefault()}function ke(e){ce(e)}function Ae(e){var n,r,i;if(!t.keyboard){e.preventDefault();return}switch(e.key){case` `:if(t.filterable)break;e.preventDefault();case`Enter`:if(!j.value?.isComposing){if(A.value){let e=M.value?.getPendingTmNode();e?Ce(e):t.filterable||(Z(),Ie())}else if(X(),t.tag&&pe.value){let e=S.value[0];if(e){let n=e[t.valueField],{value:r}=f;t.multiple&&Array.isArray(r)&&r.includes(n)||we(e)}}}e.preventDefault();break;case`ArrowUp`:if(e.preventDefault(),t.loading)return;A.value&&((n=M.value)==null||n.prev());break;case`ArrowDown`:if(e.preventDefault(),t.loading)return;A.value?(r=M.value)==null||r.next():X();break;case`Escape`:A.value&&(ie(e),Z()),(i=j.value)==null||i.focus();break}}function Ie(){var e;(e=j.value)==null||e.focus()}function Le(){var e;(e=j.value)==null||e.focusInput()}function Re(){var e;A.value&&((e=ee.value)==null||e.syncPosition())}le(),e(z(t,`options`),le);let ze={focus:()=>{var e;(e=j.value)==null||e.focus()},focusInput:()=>{var e;(e=j.value)==null||e.focusInput()},blur:()=>{var e;(e=j.value)==null||e.blur()},blurInput:()=>{var e;(e=j.value)==null||e.blurInput()}},Be=L(()=>{let{self:{menuBoxShadow:e}}=u.value;return{"--n-menu-box-shadow":e}}),Ve=c?de(`select`,void 0,Be,t):void 0;return Object.assign(Object.assign({},ze),{mergedStatus:W,mergedClsPrefix:i,mergedBordered:a,namespace:s,treeMate:D,isMounted:g(),triggerRef:j,menuRef:M,pattern:v,uncontrolledShow:k,mergedShow:A,adjustedTo:se(t),uncontrolledValue:d,mergedValue:f,followerRef:ee,localizedPlaceholder:ne,selectedOption:V,selectedOptions:B,mergedSize:re,mergedDisabled:U,focused:p,activeWithoutMenuOpen:pe,inlineThemeDisabled:c,onTriggerInputFocus:me,onTriggerInputBlur:he,handleTriggerOrMenuResize:Re,handleMenuFocus:ye,handleMenuBlur:be,handleMenuTabOut:$,handleTriggerClick:_e,handleToggle:Ce,handleDeleteOption:we,handlePatternInput:Ee,handleClear:De,handleTriggerBlur:Q,handleTriggerFocus:ve,handleKeydown:Ae,handleMenuAfterLeave:fe,handleMenuClickOutside:xe,handleMenuScroll:ke,handleMenuKeydown:Ae,handleMenuMousedown:Oe,mergedTheme:u,cssVars:c?void 0:Be,themeClass:Ve?.themeClass,onRender:Ve?.onRender})},render(){return y(`div`,{class:`${this.mergedClsPrefix}-select`},y(k,null,{default:()=>[y(ue,null,{default:()=>y(Oe,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{arrow:()=>{var e;return[(e=this.$slots).arrow?.call(e)]}})}),y(te,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===se.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{default:()=>y(ce,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{default:()=>{var e;return this.mergedShow||this.displayDirective===`show`?((e=this.onRender)==null||e.call(this),V(y(Ee,Object.assign({},this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{empty:()=>{var e;return[(e=this.$slots).empty?.call(e)]},header:()=>{var e;return[(e=this.$slots).header?.call(e)]},action:()=>{var e;return[(e=this.$slots).action?.call(e)]}}),this.displayDirective===`show`?[[oe,this.mergedShow],[Y,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[Y,this.handleMenuClickOutside,void 0,{capture:!0}]])):null}})})]}))}});export{$ as i,Me as n,Ee as r,Ie as t};